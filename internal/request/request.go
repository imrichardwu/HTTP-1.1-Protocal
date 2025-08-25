package request

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"http-1.1/internal/headers"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func (r *RequestLine) ValidHTTP() bool {
	return r.HttpVersion == "HTTP/1.1"
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	state       parserState
	Body        string
}

func getInt(headers *headers.Headers, name string, defaultvalue int) int {
	valueStr, exists := headers.Get(name)
	if !exists {
		return defaultvalue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultvalue
	}

	return value
}

var SEPERATOR = []byte("\r\n")
var ErrorMalformedRequestLine = fmt.Errorf("melform request-line")
var ErrorUnsportedHttpVersion = fmt.Errorf("unsupported http version")
var ErrorRequestInErrorState = fmt.Errorf("request in error state")

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
	}
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPERATOR)

	parts := bytes.Split(startLine, []byte(" "))

	// If we do not have method, path, and HTTP protocal
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))

	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrorMalformedRequestLine
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpParts[1]),
	}

	return rl, read, nil

}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		switch r.state {
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, ErrorRequestInErrorState
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateHeaders

		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			read += n

			if done {
				r.state = StateBody
			}

		case StateBody:
			length := getInt(r.Headers, "content-length", 0)
			if length == 0 {
				r.state = StateDone
				break outer
			}

			remaining := min(length-len(r.Body), len(currentData))
			if remaining == 0 {
				// No more data to read, but we haven't reached the expected content length
				break outer
			}

			r.Body += string(currentData[:remaining])
			read += remaining

			if len(r.Body) == length {
				r.state = StateDone
			}

		case StateDone:
			break outer
		default:
			panic("Somehow we have programmed poorly")
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	buf := make([]byte, 1024)
	bufLen := 0
	for !request.done() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			// Check if we have incomplete body content
			if err == io.EOF && request.state == StateBody {
				expectedLength := getInt(request.Headers, "content-length", 0)
				if len(request.Body) < expectedLength {
					return nil, fmt.Errorf("incomplete body: expected %d bytes, got %d", expectedLength, len(request.Body))
				}
			}
			return nil, err
		}

		bufLen += n
		readN, err := request.parse(buf[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
