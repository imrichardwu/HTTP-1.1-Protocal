package main

import (
	"fmt"
	"strings"
)

// Copy the types and function from request.go for demonstration
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var MELFORM_REQUEST_LINE = fmt.Errorf("melform request-line")
var SEPERATOR = "\r\n"

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
	return &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   parts[2],
	}, restOfMfg, nil
}

func main() {
	fmt.Println("=== parseRequestLine Function Examples ===\n")

	// Example 1: Valid GET request
	fmt.Println("Example 1: Valid GET request")
	input1 := "GET /index.html HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Mozilla/5.0\r\n\r\n"
	requestLine1, rest1, err1 := parseRequestLine(input1)

	fmt.Printf("Input: %q\n", input1)
	if err1 != nil {
		fmt.Printf("Error: %v\n", err1)
	} else if requestLine1 != nil {
		fmt.Printf("Parsed RequestLine:\n")
		fmt.Printf("  Method: %s\n", requestLine1.Method)
		fmt.Printf("  RequestTarget: %s\n", requestLine1.RequestTarget)
		fmt.Printf("  HttpVersion: %s\n", requestLine1.HttpVersion)
		fmt.Printf("Remaining data: %q\n", rest1)
	} else {
		fmt.Println("No complete request line found yet")
	}
	fmt.Println()

	// Example 2: Valid POST request
	fmt.Println("Example 2: Valid POST request")
	input2 := "POST /api/users HTTP/1.1\r\nContent-Type: application/json\r\nContent-Length: 25\r\n\r\n{\"name\":\"John\"}"
	requestLine2, rest2, err2 := parseRequestLine(input2)

	fmt.Printf("Input: %q\n", input2)
	if err2 != nil {
		fmt.Printf("Error: %v\n", err2)
	} else if requestLine2 != nil {
		fmt.Printf("Parsed RequestLine:\n")
		fmt.Printf("  Method: %s\n", requestLine2.Method)
		fmt.Printf("  RequestTarget: %s\n", requestLine2.RequestTarget)
		fmt.Printf("  HttpVersion: %s\n", requestLine2.HttpVersion)
		fmt.Printf("Remaining data: %q\n", rest2)
	} else {
		fmt.Println("No complete request line found yet")
	}
	fmt.Println()

	// Example 3: Incomplete request (no \r\n separator)
	fmt.Println("Example 3: Incomplete request (no \\r\\n separator)")
	input3 := "GET /index.html HTTP/1.1"
	requestLine3, rest3, err3 := parseRequestLine(input3)

	fmt.Printf("Input: %q\n", input3)
	if err3 != nil {
		fmt.Printf("Error: %v\n", err3)
	} else if requestLine3 != nil {
		fmt.Printf("Parsed RequestLine:\n")
		fmt.Printf("  Method: %s\n", requestLine3.Method)
		fmt.Printf("  RequestTarget: %s\n", requestLine3.RequestTarget)
		fmt.Printf("  HttpVersion: %s\n", requestLine3.HttpVersion)
		fmt.Printf("Remaining data: %q\n", rest3)
	} else {
		fmt.Println("No complete request line found yet - need more data")
		fmt.Printf("Buffered data: %q\n", rest3)
	}
	fmt.Println()

	// Example 4: Malformed request (too few parts)
	fmt.Println("Example 4: Malformed request (too few parts)")
	input4 := "GET /index.html\r\nHost: example.com\r\n"
	requestLine4, rest4, err4 := parseRequestLine(input4)

	fmt.Printf("Input: %q\n", input4)
	if err4 != nil {
		fmt.Printf("Error: %v\n", err4)
		fmt.Printf("Remaining data: %q\n", rest4)
	} else if requestLine4 != nil {
		fmt.Printf("Parsed RequestLine:\n")
		fmt.Printf("  Method: %s\n", requestLine4.Method)
		fmt.Printf("  RequestTarget: %s\n", requestLine4.RequestTarget)
		fmt.Printf("  HttpVersion: %s\n", requestLine4.HttpVersion)
		fmt.Printf("Remaining data: %q\n", rest4)
	} else {
		fmt.Println("No complete request line found yet")
	}
	fmt.Println()

	// Example 5: Malformed request (too many parts)
	fmt.Println("Example 5: Malformed request (too many parts)")
	input5 := "GET /index.html HTTP/1.1 extra\r\nHost: example.com\r\n"
	requestLine5, rest5, err5 := parseRequestLine(input5)

	fmt.Printf("Input: %q\n", input5)
	if err5 != nil {
		fmt.Printf("Error: %v\n", err5)
		fmt.Printf("Remaining data: %q\n", rest5)
	} else if requestLine5 != nil {
		fmt.Printf("Parsed RequestLine:\n")
		fmt.Printf("  Method: %s\n", requestLine5.Method)
		fmt.Printf("  RequestTarget: %s\n", requestLine5.RequestTarget)
		fmt.Printf("  HttpVersion: %s\n", requestLine5.HttpVersion)
		fmt.Printf("Remaining data: %q\n", rest5)
	} else {
		fmt.Println("No complete request line found yet")
	}
	fmt.Println()

	// Example 6: DELETE request with query parameters
	fmt.Println("Example 6: DELETE request with query parameters")
	input6 := "DELETE /api/users/123?force=true HTTP/1.1\r\nAuthorization: Bearer token123\r\n\r\n"
	requestLine6, rest6, err6 := parseRequestLine(input6)

	fmt.Printf("Input: %q\n", input6)
	if err6 != nil {
		fmt.Printf("Error: %v\n", err6)
	} else if requestLine6 != nil {
		fmt.Printf("Parsed RequestLine:\n")
		fmt.Printf("  Method: %s\n", requestLine6.Method)
		fmt.Printf("  RequestTarget: %s\n", requestLine6.RequestTarget)
		fmt.Printf("  HttpVersion: %s\n", requestLine6.HttpVersion)
		fmt.Printf("Remaining data: %q\n", rest6)
	} else {
		fmt.Println("No complete request line found yet")
	}
}
