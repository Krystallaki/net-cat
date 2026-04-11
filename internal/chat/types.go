// Package chat implements the chat logic for TCPChat.
package chat

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
