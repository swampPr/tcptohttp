// Package main provides main  INFO:  tcp_client
package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	ln, err := net.Dial("tcp", ":42069")
	if err != nil {
		log.Fatalf("ERROR DIALING PORT: %v", err)
	}

	defer func() { _ = ln.Close() }()

	f, err := os.ReadFile("../messages.txt")
	if err != nil {
		log.Fatalf("ERROR READING FILE: %v", err)
	}

	_, _ = fmt.Fprintf(ln, "%v", f)
}
