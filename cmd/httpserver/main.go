package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"http-1.1/internal/headers"
	"http-1.1/internal/request"
	"http-1.1/internal/response"
	"http-1.1/internal/server"
)

const port = 42069

func respond400() []byte {
	return []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
}

func respond500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}

func respond200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}

func toStr(bytes []byte) string {
	out := ""
	for _, b := range bytes {
		out += fmt.Sprintf("%02x", b)
	}
	return out
}

func main() {
	s, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		h := response.GetDefaultHeaders(0)
		h.Set("Content-Type", "text/html")

		body := respond200()
		status := response.StatusOk

		if req.RequestLine.RequestTarget == "/yourproblem" {
			body = respond400()
			status = response.StatusBadRequest
		} else if req.RequestLine.Method == "/myproblem" {
			body = respond500()
			status = response.StatusInternalServerError
		} else if req.RequestLine.RequestTarget == "/video" {
			f, _ := os.ReadFile("assets/vim.mp4")
			h.Replace("Content-Type", "video/mp4")
			h.Replace("Content-Length", fmt.Sprintf("%d", len(f)))

			w.WriteStatusLine(response.StatusOk)
			w.WriteHeaders(*h)
			w.WriteBody(f)

		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			target := req.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
			if err != nil {
				body = respond500()
				status = response.StatusInternalServerError
			} else {
				h.Delete("Content-Length")
				h.Set("transfer-encoding", "chunked")
				h.Replace("Content-Type", "text/plain")
				h.Set("Trailer", "X-Content-SHA256")
				h.Set("Trailer", "X-Content-Length")
				w.WriteHeaders(*h)

				fullbody := []byte{}
				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break
					}
					fullbody = append(fullbody, data[:n]...)
					w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n"))
				trailer := headers.NewHeaders()
				out := sha256.Sum256(fullbody)
				trailer.Set("X-Content-SHA256", toStr(out[:]))
				trailer.Set("X-Content-Length", fmt.Sprintf("%d", len(fullbody)))
				w.WriteHeaders(*trailer)
				return
			}
		}

		h.Replace("Content-Length", fmt.Sprintf("%d", len(body)))
		h.Replace("Content-Type", "text/html")
		w.WriteStatusLine(response.StatusCode(status))
		w.WriteHeaders(*h)
		w.WriteBody(body)
	})

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
