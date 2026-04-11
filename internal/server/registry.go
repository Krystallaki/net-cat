package server

import (
	"net-cat/internal/chat"
	"sync"
)

// registry is a thread-safe set of connected clients.
type registry struct {
	mu      sync.Mutex
	clients map[*chat.Client]struct{}
}

// TODO (Vasiliki): Add(c *chat.Client)    — register a new client
// TODO (Vasiliki): Remove(c *chat.Client) — deregister a client on disconnect
// TODO (Vasiliki): All() []*chat.Client   — return a snapshot of all clients (used by Broadcast)
