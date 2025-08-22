package request

import (
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ERROR_BAD_START_LINE = fmt.Printf("bad start line")

func RequestFromReader(reader io.Reader) (*Request, error) {

}
