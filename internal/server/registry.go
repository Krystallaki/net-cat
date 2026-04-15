package server

import (
	"net-cat/internal/client"
	"sync"
)

type registry struct {
	mu      sync.Mutex
	clients map[*client.Client]struct{}
}
func (r *registry) Add(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[c] = struct{}{}
}
func (r *registry) Remove(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, c)
}
func (r *registry) All() (s []*client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.clients {
		s = append(s, c)
	}
	return s
}
