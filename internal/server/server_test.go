package server

import (
	"net"
	"testing"
	"time"
)

// startServer launches a Server on the given port in a background goroutine
// and returns a dial helper that connects to it.
func startServer(t *testing.T, port string) func() net.Conn {
	t.Helper()
	s := &Server{}
	go s.Start(port)
	time.Sleep(20 * time.Millisecond)
	return func() net.Conn {
		conn, err := net.Dial("tcp", ":"+port)
		if err != nil {
			t.Fatalf("failed to dial server: %v", err)
		}
		return conn
	}
}

func TestStart_ListenerBinds(t *testing.T) {
	dial := startServer(t, "9001")
	conn := dial()
	defer conn.Close()
}

func TestStart_AcceptsUpToTenClients(t *testing.T) {
	dial := startServer(t, "9002")
	conns := make([]net.Conn, 10)
	for i := range conns {
		conns[i] = dial()
		defer conns[i].Close()
	}
}

func TestStart_RejectsEleventhClient(t *testing.T) {
	dial := startServer(t, "9003")
	conns := make([]net.Conn, 10)
	for i := range conns {
		conns[i] = dial()
		defer conns[i].Close()
	}

	extra := dial()
	defer extra.Close()

	buf := make([]byte, 64)
	extra.SetReadDeadline(time.Now().Add(time.Second))
	n, _ := extra.Read(buf)

	if got := string(buf[:n]); got != "Chat is full. Try again later.\n" {
		t.Errorf("expected rejection message, got %q", got)
	}
}
