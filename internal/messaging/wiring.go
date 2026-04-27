// Package messaging handles message formatting, history, delivery, and clean
// goroutine teardown for TCPChat.
package messaging

import (
	"bufio"
	"net-cat/internal/client"
	"time"
)

// RunMessageLoop reads lines from c until the connection closes or an error occurs.
// all is a snapshot of connected clients used for broadcast.
// On each message it formats, records in hist, and broadcasts to all other clients.
// The loop exits cleanly when the scanner stops — either on disconnect or a read error —
// preventing any write to a closed net.Conn.
func RunMessageLoop(c *client.Client, all func() []*client.Client, hist *History) {
	scanner := bufio.NewScanner(c.Conn)
	for scanner.Scan() {
		body := scanner.Text()
		if IsEmpty(body) {
			continue
		}
		msg := FormatMessage(time.Now(), c.Name, body)
		hist.Add(msg)
		client.Broadcast(all(), msg+"\n", c)
	}
}
