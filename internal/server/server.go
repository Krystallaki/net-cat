// Package server implements the TCP listener and client lifecycle for TCPChat.
package server

// Server manages the TCP listener and connected clients.
type Server struct {
	reg registry
}

// TODO (Vasiliki): Start(port string) error — bind listener, accept connections, enforce max 10
