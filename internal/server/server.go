package server

import (
	"bytes"
	"fmt"
	"io"
	"net"

	"http-1.1/internal/request"
	"http-1.1/internal/response"
)

type Server struct {
	closed  bool
	handler Handler
}

type HandlerError struct {
	StatusCode response.StatusCode
	Message    string
}
type Handler func(w io.Writer, req *request.Request) *HandlerError

func runConnection(_s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	headers := response.GetDefaultHeaders(0)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		response.WriteStatusLine(conn, response.StatusBadRequest)
		response.WriteHeaders(conn, headers)
		return
	}
	writer := bytes.NewBuffer([]byte{})
	HandlerError := _s.handler(writer, r)

	var body []byte
	var status response.StatusCode = response.StatusOk
	if HandlerError != nil {
		status = HandlerError.StatusCode
		body = []byte(HandlerError.Message)
	} else {
		body = writer.Bytes()
	}

	headers.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
	response.WriteStatusLine(conn, status)
	response.WriteHeaders(conn, headers)
	conn.Write(body)

}

func runServer(s *Server, listener net.Listener) {

	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}

		if err != nil {
			return
		}
		go runConnection(s, conn)
	}

}

func Serve(port uint, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{closed: false, handler: handler}
	go runServer(server, listener)
	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
