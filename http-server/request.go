package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
	Params  map[string]string
	Query   map[string]string
}

func ParseRequest(reader *bufio.Reader) (HTTPRequest, error) {
	message, err := reader.ReadString('\n')
	if err != nil {
		return HTTPRequest{}, fmt.Errorf("read request line: %w", err)
	}

	requestLine := strings.TrimSpace(message)
	parts := strings.Split(requestLine, " ")
	if len(parts) < 2 {
		return HTTPRequest{}, fmt.Errorf("malformed request line: %s", requestLine)
	}

	rawPath := parts[1]
	query := make(map[string]string)

	if idx := strings.Index(rawPath, "?"); idx != -1 {
		queryString := rawPath[idx+1:]
		rawPath = rawPath[:idx]

		for _, pair := range strings.Split(queryString, "&") {
			if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 {
				query[kv[0]] = kv[1]
			}
		}
	}

	req := HTTPRequest{
		Method:  parts[0],
		Path:    rawPath,
		Headers: make(map[string]string),
		Params:  make(map[string]string),
		Query:   query,
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return req, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		if idx := strings.Index(line, ":"); idx != -1 {
			key := strings.TrimSpace(line[:idx])
			val := strings.TrimSpace(line[idx+1:])
			req.Headers[key] = val
		}
	}

	if contentLength, ok := req.Headers["Content-Length"]; ok {
		bodySize, err := strconv.Atoi(contentLength)
		if err != nil {
			return req, fmt.Errorf("invalid Content-Length: %w", err)
		}
		body := make([]byte, bodySize)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			return req, fmt.Errorf("read body: %w", err)
		}
		req.Body = string(body)
	}

	return req, nil
}
