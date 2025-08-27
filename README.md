# HTTP/1.1 Server Implementation

A custom HTTP/1.1 server implementation written in Go that demonstrates low-level HTTP protocol parsing and handling.

## Overview

This project implements a complete HTTP/1.1 server from scratch, including:

- Raw TCP connection handling
- HTTP request parsing (request line, headers, body)
- HTTP response generation
- Chunked transfer encoding support
- Static file serving
- HTTP proxy functionality

## Features

- **Full HTTP/1.1 Protocol Support**: Parses HTTP requests according to RFC specifications
- **Concurrent Request Handling**: Uses goroutines to handle multiple simultaneous connections
- **Chunked Transfer Encoding**: Supports streaming responses with trailer headers
- **Static File Serving**: Serves video files (MP4) with proper content headers
- **HTTP Proxy**: Proxies requests to httpbin.org with chunked encoding
- **Comprehensive Testing**: Unit tests for request parsing and header handling
- **Graceful Shutdown**: Handles SIGINT/SIGTERM for clean server shutdown

## Project Structure

```
.
├── cmd/
│   ├── httpserver/          # Main HTTP server application
│   │   └── main.go
│   └── tcplistener/         # TCP connection listener for debugging
│       └── main.go
├── internal/
│   ├── headers/             # HTTP header parsing and management
│   │   ├── headers.go
│   │   └── headers_test.go
│   ├── request/             # HTTP request parsing
│   │   ├── request.go
│   │   └── request_test.go
│   ├── response/            # HTTP response generation
│   │   └── response.go
│   └── server/              # Server implementation
│       └── server.go
├── assets/
│   └── vim.mp4             # Example video file for serving
├── go.mod
├── go.sum
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24.4 or later
- Make sure port 42069 is available

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd HTTP-1.1
```

2. Install dependencies:

```bash
go mod download
```

### Running the Server

Start the HTTP server:

```bash
go run cmd/httpserver/main.go
```

The server will start on port `42069` and log its status.

For debugging HTTP requests, you can also run the TCP listener:

```bash
go run cmd/tcplistener/main.go
```

### Testing the Server

Once the server is running, you can test it with curl:

#### Basic requests:

```bash
# Success response
curl http://localhost:42069/

# 400 Bad Request
curl http://localhost:42069/yourproblem

# 500 Internal Server Error
curl http://localhost:42069/myproblem
```

#### Video file serving:

```bash
# Download the video file
curl http://localhost:42069/video -o output.mp4
```

#### HTTP proxy (chunked encoding):

```bash
# Proxy request to httpbin.org
curl http://localhost:42069/httpbin/get
curl http://localhost:42069/httpbin/headers
```

### Running Tests

Execute the test suite:

```bash
go test ./...
```

Run tests with verbose output:

```bash
go test -v ./...
```

## API Endpoints

| Endpoint       | Method | Description                                           |
| -------------- | ------ | ----------------------------------------------------- |
| `/`            | GET    | Returns a success HTML page                           |
| `/yourproblem` | GET    | Returns a 400 Bad Request error page                  |
| `/myproblem`   | Any    | Returns a 500 Internal Server Error page              |
| `/video`       | GET    | Serves the vim.mp4 file                               |
| `/httpbin/*`   | Any    | Proxies requests to httpbin.org with chunked encoding |

## Architecture

### Request Processing Flow

1. **TCP Connection**: Server accepts incoming TCP connections on port 42069
2. **Request Parsing**: Raw bytes are parsed into structured HTTP request objects
3. **Request Routing**: Based on the request target, appropriate handler is invoked
4. **Response Generation**: HTTP response is constructed with proper headers and body
5. **Connection Cleanup**: Connection is closed after response is sent

### Key Components

- **`server.Server`**: Main server struct that manages TCP connections and request routing
- **`request.Request`**: Represents a parsed HTTP request with request line, headers, and body
- **`response.Writer`**: Handles HTTP response writing with status lines, headers, and body
- **`headers.Headers`**: Manages HTTP header parsing and manipulation

### HTTP Features Implemented

- **Request Line Parsing**: Method, URI, and HTTP version extraction
- **Header Parsing**: Case-insensitive header names with value concatenation for duplicates
- **Content-Length**: Proper body parsing based on Content-Length header
- **Transfer-Encoding**: Chunked encoding support for streaming responses
- **Trailer Headers**: Support for trailer headers in chunked responses
- **Content-Type Detection**: Automatic content type setting based on file extensions

## Implementation Details

### Error Handling

The server implements proper HTTP error responses:

- Malformed requests return 400 Bad Request
- Server errors return 500 Internal Server Error
- Missing resources could be extended to return 404 Not Found

### Concurrency

Each incoming connection is handled in a separate goroutine, allowing the server to process multiple requests simultaneously without blocking.

### Memory Management

The server uses buffered reading to efficiently parse incoming requests while managing memory usage for large request bodies.

## License

This project is for educational purposes, demonstrating HTTP/1.1 protocol implementation in Go.

## Acknowledgments

- HTTP/1.1 specification (RFC 7230-7237)
- Go standard library for networking primitives
- httpbin.org for testing HTTP proxy functionality
