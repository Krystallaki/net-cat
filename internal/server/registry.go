// Package server implements the TCP chat server, managing client connections,
// message broadcasting, and the lifecycle of each connected session.
package server

import (
	"net-cat/internal/client"
	"sync"
)

// registry is a thread-safe set of connected clients.
type registry struct {
	mu      sync.Mutex
	clients map[*client.Client]struct{}
}

// Add registers c in the registry.
func (r *registry) Add(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

// Remove deregisters c from the registry.
func (r *registry) Remove(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

// All returns a snapshot of all currently registered clients.
func (r *registry) All() (s []*client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.clients {
		s = append(s, c)
	}
	return s
}
