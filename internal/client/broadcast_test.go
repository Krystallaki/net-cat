//nolint:errcheck
package client

import (
	"net"
	"testing"
	"time"
)

func TestBroadcast_SendsToIncluded(t *testing.T) {
	receiverServer, receiverClient := net.Pipe()
	excludedServer, excludedClient := net.Pipe()
	defer receiverServer.Close()
	defer receiverClient.Close()
	defer excludedServer.Close()
	defer excludedClient.Close()

	receiver := &Client{Conn: receiverServer, Name: "receiver"}
	excluded := &Client{Conn: excludedServer, Name: "excluded"}

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 64)
		n, _ := receiverClient.Read(buf)
		done <- string(buf[:n])
	}()

	Broadcast([]*Client{receiver, excluded}, "hello\n", excluded)

	got := <-done
	if got != "hello\n" {
		t.Errorf("expected 'hello\\n', got %q", got)
	}
}

func TestBroadcast_ExcludedDoesNotReceive(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer serverSide.Close()
	defer clientSide.Close()

	excluded := &Client{Conn: serverSide, Name: "excluded"}
	Broadcast([]*Client{excluded}, "hello\n", excluded)

	clientSide.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	buf := make([]byte, 64)
	n, _ := clientSide.Read(buf)
	if n > 0 {
		t.Errorf("excluded client received %q", string(buf[:n]))
	}
}

func TestBroadcast_NilExcludeSendsToAll(t *testing.T) {
	server1, client1 := net.Pipe()
	server2, client2 := net.Pipe()
	defer server1.Close()
	defer client1.Close()
	defer server2.Close()
	defer client2.Close()

	c1 := &Client{Conn: server1, Name: "user1"}
	c2 := &Client{Conn: server2, Name: "user2"}

	done := make(chan string, 2)
	go func() {
		buf := make([]byte, 64)
		n, _ := client1.Read(buf)
		done <- string(buf[:n])
	}()
	go func() {
		buf := make([]byte, 64)
		n, _ := client2.Read(buf)
		done <- string(buf[:n])
	}()

	Broadcast([]*Client{c1, c2}, "hello\n", nil)

	for range 2 {
		if got := <-done; got != "hello\n" {
			t.Errorf("expected 'hello\\n', got %q", got)
		}
	}
}
