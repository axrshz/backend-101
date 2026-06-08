package main

import "log"

type Middleware func(next func(HTTPRequest) HTTPResponse) func(HTTPRequest) HTTPResponse

var middlewares []Middleware

func Use(m Middleware) {
	middlewares = append(middlewares, m)
}

func applyMiddleware(handler func(HTTPRequest) HTTPResponse) func(HTTPRequest) HTTPResponse {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func LoggingMiddleware(next func(HTTPRequest) HTTPResponse) func(HTTPRequest) HTTPResponse {
	return func(req HTTPRequest) HTTPResponse {
		resp := next(req)
		log.Printf("[%s] %s -> %d", req.Method, req.Path, resp.StatusCode)
		return resp
	}
}

func AuthMiddleware(next func(HTTPRequest) HTTPResponse) func(HTTPRequest) HTTPResponse {
	return func(req HTTPRequest) HTTPResponse {
		if _, ok := req.Headers["Authorization"]; !ok {
			return HTTPResponse{
				StatusCode: 401,
				Headers:    map[string]string{"Content-Type": "text/plain"},
				Body:       "401 Unauthorized\n",
			}
		}
		return next(req)
	}
}
