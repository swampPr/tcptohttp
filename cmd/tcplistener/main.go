// Package main provides main  INFO:  TCP-TO-HTTP (https://www.youtube.com/watch?v=FknTw9bJsXM)
package main

import (
	"fmt"
	"log"
	"net"

	"tcp_http/internal/request"
)

func main() {
	//#nosec
	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("[ERROR OPENING PORT]: %v", err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("ERROR ACCEPTING CONNECTION: %v", err)
		}

		rl, err := request.RequestFromReader(conn)
		if err != nil {
			_, _ = fmt.Fprintf(conn, "ERROR: Could not parse HTTP reqest: %s", err)
			log.Printf("ERROR: Could not parse HTTP request: %s", err)
		}

		fmt.Printf(`
		Request Line:
		- Method: %s
		- Target: %s
		- Version: %s
		`, rl.RequestLine.Method, rl.RequestLine.RequestTarget, rl.RequestLine.HTTPVersion)

	}
}
