package main

import "strings"

type Route struct {
	Method  string
	Pattern string
	Handler func(HTTPRequest) HTTPResponse
}

var routes = []Route{
	{"GET", "/hello", helloHandler},
	{"GET", "/users", getUsersHandler},
	{"GET", "/users/:id", getUserHandler},
	{"POST", "/users", createUserHandler},
}

func routeRequest(req HTTPRequest) HTTPResponse {
	for _, route := range routes {
		if route.Method != req.Method {
			continue
		}

		params, ok := matchPath(route.Pattern, req.Path)
		if ok {
			req.Params = params
			handler := applyMiddleware(route.Handler)
			return handler(req)
		}
	}

	return HTTPResponse{
		StatusCode: 404,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       "404 Not Found\n",
	}
}

func matchPath(pattern, path string) (map[string]string, bool) {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return nil, false
	}

	params := make(map[string]string)
	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			params[part[1:]] = pathParts[i]
		} else if part != pathParts[i] {
			return nil, false
		}
	}

	return params, true
}
