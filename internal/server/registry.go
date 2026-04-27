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
// It is safe to call concurrently from multiple goroutines.
func (r *registry) Add(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}

// Remove deregisters c from the registry.
// If c is not present, Remove is a no-op.
func (r *registry) Remove(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}

// All returns a snapshot of all currently registered clients.
// The returned slice is independent of the internal map — callers may iterate
// it safely without holding the lock.
func (r *registry) All() (s []*client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.clients {
		s = append(s, c)
	}
	return s
}

// IsFull reports whether the registry has reached the maximum of 10 clients.
// The check is performed atomically under the lock to prevent a race condition
// between reading the count and registering a new client.
func (r *registry) IsFull() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.clients) >= 10
}
