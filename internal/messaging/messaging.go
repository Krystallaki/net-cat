// Package messaging handles message formatting, history, and delivery for TCPChat.
package messaging

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// FormatMessage returns a chat message string in the format [timestamp][name]:body.
func FormatMessage(ts time.Time, name, body string) string {
	return fmt.Sprintf("[%s][%s]:%s", ts.Format("2006-01-02 15:04:05"), name, body)
}

// IsEmpty reports whether a message body is empty or whitespace-only.
func IsEmpty(body string) bool {
	return strings.TrimSpace(body) == ""
}

// History is a thread-safe ordered log of formatted message strings.
type History struct {
	mu       sync.Mutex
	messages []string
}

// Add appends a formatted message string to the history.
func (h *History) Add(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.messages = append(h.messages, msg)
}

// Replay sends all history messages to conn, each followed by a newline.
func (h *History) Replay(conn net.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, msg := range h.messages {
		fmt.Fprintln(conn, msg) //nolint:errcheck
	}
}
