package main

import (
	"fmt"
	"net"
)

type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func (r HTTPResponse) WriteTo(conn net.Conn) error {
	statusText := map[int]string{
		200: "OK",
		404: "Not Found",
		500: "Internal Server Error",
	}

	text, ok := statusText[r.StatusCode]
	if !ok {
		text = "Unknown"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, text)

	headers := ""
	for key, val := range r.Headers {
		headers += fmt.Sprintf("%s: %s\r\n", key, val)
	}
	headers += fmt.Sprintf("Content-Length: %d\r\n", len(r.Body))
	headers += "Connection: close\r\n"

	raw := statusLine + headers + "\r\n" + r.Body
	_, err := conn.Write([]byte(raw))
	return err
}
