// Package client defines the shared types and lifecycle logic for TCPChat.
package client

import "net"

// Client represents a connected chat client.
type Client struct {
	Conn net.Conn
	Name string
}

// Message represents a formatted chat message.
type Message struct {
	Timestamp string
	Sender    string
	Content   string
}
