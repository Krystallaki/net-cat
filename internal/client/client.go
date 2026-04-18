// Package client defines the shared types and lifecycle logic for TCPChat.
package client

import (
	"net"
	"time"
)

// Client represents a connected chat client.
type Client struct {
	Conn net.Conn
	Name string
}

// Message represents a formatted chat message.
// Timestamp is stored as time.Time and formatted at display time.
type Message struct {
	Timestamp time.Time
	Sender    string
	Body      string
}
