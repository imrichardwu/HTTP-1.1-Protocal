package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
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
}

var SEPERATOR = "\r\n"
var MELFORM_REQUEST_LINE = fmt.Errorf("melform request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPERATOR)
	if idx == -1 {
		return nil, b, nil
	}
	startLine := b[:idx]
	restOfMfg := b[idx+len(SEPERATOR):]

	parts := strings.Split(startLine, " ")

	// If we do not have method, path, and HTTP protocal
	if len(parts) != 3 {
		return nil, restOfMfg, MELFORM_REQUEST_LINE
	}
	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}

	if !rl.ValidHTTP() {
		return nil, restOfMfg, MELFORM_REQUEST_LINE
	}

	return rl, restOfMfg, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("unable to io.ReadALL: %w", err))
	}

	str := string(data)
	rl, _, err := parseRequestLine(str)

	return &Request{RequestLine: *rl}, err
}
