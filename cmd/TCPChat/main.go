// Package main is the entry point for the TCPChat server.
package main

import (
	"fmt"
	"os"
)

func main() {
	// TODO (Vasiliki): parse args, validate port, start server
	// Usage: go run . [port]
	// Default port: 8989
	// Error if more than 1 arg: [USAGE]: ./TCPChat $port
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
}
