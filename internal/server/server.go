package server

import "sync"

// Server holds the TCP listener state and connected clients.
type Server struct {
	clients map[*Client]struct{}
	history []string
	mu      sync.Mutex
}

// TODO (Vasiliki): implement Start(port string) error — bind listener, accept connections, enforce max 10
