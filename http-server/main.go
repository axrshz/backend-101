package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	Use(LoggingMiddleware)
	Use(AuthMiddleware)

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	fmt.Println("Server listening on :8090")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	req, err := ParseRequest(reader)
	if err != nil {
		log.Printf("Parse error: %v", err)
		return
	}

	fmt.Printf("[%s] %s\n", req.Method, req.Path)

	resp := routeRequest(req)
	err = resp.WriteTo(conn)
	if err != nil {
		log.Printf("Write error: %v", err)
	}
}
