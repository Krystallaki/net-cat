package client

import "fmt"

// Broadcast sends msg to every client in clients except exclude.
// Pass exclude as nil to send to all clients (e.g. server announcements).
func Broadcast(clients []*Client, msg string, exclude *Client) {
	for _, c := range clients {
		if c == exclude {
			continue
		}
		fmt.Fprint(c.Conn, msg) //nolint:errcheck
	}
}
