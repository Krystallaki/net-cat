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

// TODO (Vasiliki): Add(c *client.Client)    — register a new client
// TODO (Vasiliki): Remove(c *client.Client) — deregister a client on disconnect
// TODO (Vasiliki): All() []*client.Client   — return a snapshot of all clients (used by Broadcast)
func(r *registry) Add(c *client.Client) {
r.mu.Lock()
defer r.mu.Unlock()
r.clients[c] = struct{}{}
}