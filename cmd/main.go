// Package main is the entry point for the TCPChat server.
// It parses command-line arguments, resolves the port to listen on,
// and starts the TCP server. If more than one argument is provided,
// it prints a usage message and exits.
package main

import (
	"fmt"
	"net-cat/internal/server"
	"os"
	"strconv"
)

const defaultPort = "8989"

// validatePort reports whether port is a valid TCP port number (1–65535).
func validatePort(port string) bool {
	n, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return n >= 1 && n <= 65535
}

// resolvePort returns defaultPort when no argument is provided, or the given argument otherwise.
func resolvePort(args []string) string {
	if len(args) == 0 {
		return defaultPort
	}
	return args[0]
}

// parseArgs validates os.Args and returns the port to listen on.
// It returns "", false if more than one argument is provided or the port is invalid.
func parseArgs(args []string) (string, bool) {
	if len(args) > 1 {
		return "", false
	}
	port := resolvePort(args)
	if !validatePort(port) {
		return "", false
	}
	return port, true
}

// main parses arguments, resolves the port, and starts the TCP server.
// It exits with a non-zero status if the arguments are invalid or the server fails to start.
func main() {
	port, ok := parseArgs(os.Args[1:])
	if !ok {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	srv := &server.Server{}
	err := srv.Start(port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
