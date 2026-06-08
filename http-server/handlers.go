package main

import "fmt"

func helloHandler(req HTTPRequest) HTTPResponse {
	return HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       "Hello, World!\n",
	}
}

func getUsersHandler(req HTTPRequest) HTTPResponse {
	users := `[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]`

	if name, ok := req.Query["name"]; ok {
		users = fmt.Sprintf(`[{"id":1,"name":"%s"}]`, name)
	}

	return HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       users + "\n",
	}
}

func getUserHandler(req HTTPRequest) HTTPResponse {
	id := req.Params["id"]
	return HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf(`{"id":%s,"name":"User %s"}`+"\n", id, id),
	}
}

func createUserHandler(req HTTPRequest) HTTPResponse {
	return HTTPResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf(`{"message":"Created","data":%s}`+"\n", req.Body),
	}
}
