// Package messaging handles message formatting, history, and delivery for TCPChat.
package messaging

import (
	"fmt"
	"time"
)

// FormatMessage returns a chat message string in the format [timestamp][name]:body.
func FormatMessage(ts time.Time, name, body string) string {
	return fmt.Sprintf("[%s][%s]:%s", ts.Format("2006-01-02 15:04:05"), name, body)
}
