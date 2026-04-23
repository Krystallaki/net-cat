package messaging

import (
	"bufio"
	"net-cat/internal/client"
	"sync"
	"time"
)

// Registry is a thread-safe set of connected clients owned by the messaging layer.
type Registry struct {
	mu      sync.Mutex
	clients []*client.Client
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Add registers c in the registry.
func (r *Registry) Add(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients = append(r.clients, c)
}

// Remove deregisters c from the registry.
func (r *Registry) Remove(c *client.Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	updated := make([]*client.Client, 0, len(r.clients))
	for _, existing := range r.clients {
		if existing != c {
			updated = append(updated, existing)
		}
	}
	r.clients = updated
}

// All returns a snapshot of all currently registered clients.
func (r *Registry) All() []*client.Client {
	r.mu.Lock()
	defer r.mu.Unlock()
	snapshot := make([]*client.Client, len(r.clients))
	copy(snapshot, r.clients)
	return snapshot
}

// HandleClient runs the full lifecycle for a single connected client:
// welcome banner → name → history replay → join notify → message loop → leave notify → cleanup.
func HandleClient(c *client.Client, reg *Registry, hist *History) {
	defer reg.Remove(c)
	defer c.Conn.Close()

	client.SendWelcome(c.Conn)

	name, err := client.ReadName(c.Conn)
	if err != nil {
		return
	}
	c.Name = name

	hist.Replay(c.Conn)
	client.NotifyJoin(reg.All(), c)

	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		body := scanner.Text()
		if IsEmpty(body) {
			continue
		}
		msg := FormatMessage(time.Now(), c.Name, body)
		hist.Add(msg)
		client.Broadcast(reg.All(), msg+"\n", c)
	}

	client.NotifyLeave(reg.All(), c)
}
