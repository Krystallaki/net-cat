package server

import (
	"bufio"
	"net"
	"strings"
	"testing"
)

// net.Pipe() δημιουργεί δύο συνδεδεμένα net.Conn.
// Ό,τι γράψουμε στο serverSide φαίνεται στο clientSide.

func TestSendWelcome_ContainsBanner(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer serverSide.Close()
	defer clientSide.Close()

	go func() {
		sendWelcome(serverSide)
		serverSide.Close()
	}()

	var received strings.Builder
	scanner := bufio.NewScanner(clientSide)
	for scanner.Scan() {
		line := scanner.Text()
		received.WriteString(line + "\n")
		if strings.Contains(line, "[ENTER YOUR NAME]:") {
			break
		}
	}

	got := received.String()

	if !strings.Contains(got, "Welcome to TCP-Chat!") {
		t.Errorf("expected 'Welcome to TCP-Chat!' in output, got:\n%s", got)
	}
	if !strings.Contains(got, "[ENTER YOUR NAME]:") {
		t.Errorf("expected '[ENTER YOUR NAME]:' in output, got:\n%s", got)
	}
}
