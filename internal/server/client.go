// Package server implements the TCPChat server logic.
package server

import "net"

// Client represents a connected chat client.
type Client struct {
	Conn net.Conn
	Name string
}
