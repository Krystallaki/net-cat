package messaging

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"net-cat/internal/client"
)

func TestFormatMessage(t *testing.T) {
	ts := time.Date(2020, 1, 20, 15, 48, 41, 0, time.UTC)
	got := FormatMessage(ts, "Yenlik", "hello")
	want := "[2020-01-20 15:48:41][Yenlik]:hello"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIsEmpty_EmptyString(t *testing.T) {
	if !IsEmpty("") {
		t.Error("expected empty string to be empty")
	}
}

func TestIsEmpty_WhitespaceOnly(t *testing.T) {
	if !IsEmpty("   ") {
		t.Error("expected whitespace-only string to be empty")
	}
}

func TestIsEmpty_NonEmpty(t *testing.T) {
	if IsEmpty("hello") {
		t.Error("expected non-empty string to not be empty")
	}
}

func TestHistory_AppendAndReplay(t *testing.T) {
	h := &History{}
	h.Add("[2020-01-20 15:48:41][Yenlik]:hello")
	h.Add("[2020-01-20 15:48:42][Yenlik]:world")

	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	done := make(chan string, 1)
	go func() {
		var sb strings.Builder
		scanner := bufio.NewScanner(clientSide)
		for scanner.Scan() {
			sb.WriteString(scanner.Text() + "\n")
		}
		done <- sb.String()
	}()

	h.Replay(serverSide)
	serverSide.Close()

	got := <-done
	want := "[2020-01-20 15:48:41][Yenlik]:hello\n[2020-01-20 15:48:42][Yenlik]:world\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHandleClient_SetsNameAndBroadcastsJoin(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	reg := NewRegistry()
	hist := &History{}

	c := &client.Client{Conn: serverSide}
	reg.Add(c)

	done := make(chan struct{})
	go func() {
		HandleClient(c, reg, hist)
		close(done)
	}()

	// read and discard the welcome banner
	buf := make([]byte, 512)
	clientSide.Read(buf)

	// send name
	fmt.Fprint(clientSide, "Theo\n")

	// read join notification sent to all (only one client so nobody else receives it,
	// but we can assert the name was set)
	time.Sleep(20 * time.Millisecond)
	if c.Name != "Theo" {
		t.Errorf("expected client name 'Theo', got %q", c.Name)
	}

	clientSide.Close()
	<-done
}

func TestHandleClient_EmptyMessageNotBroadcast(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	reg := NewRegistry()
	hist := &History{}

	c := &client.Client{Conn: serverSide}
	reg.Add(c)

	done := make(chan struct{})
	go func() {
		HandleClient(c, reg, hist)
		close(done)
	}()

	buf := make([]byte, 512)
	clientSide.Read(buf)
	fmt.Fprint(clientSide, "Theo\n")
	time.Sleep(20 * time.Millisecond)

	// send empty message
	fmt.Fprint(clientSide, "\n")
	time.Sleep(20 * time.Millisecond)

	hist.mu.Lock()
	count := len(hist.messages)
	hist.mu.Unlock()

	if count != 0 {
		t.Errorf("expected 0 messages in history after empty send, got %d", count)
	}

	clientSide.Close()
	<-done
}

func TestHandleClient_MessageAddedToHistory(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer clientSide.Close()

	reg := NewRegistry()
	hist := &History{}

	c := &client.Client{Conn: serverSide}
	reg.Add(c)

	done := make(chan struct{})
	go func() {
		HandleClient(c, reg, hist)
		close(done)
	}()

	buf := make([]byte, 512)
	clientSide.Read(buf)
	fmt.Fprint(clientSide, "Theo\n")
	time.Sleep(20 * time.Millisecond)

	fmt.Fprint(clientSide, "hello world\n")
	time.Sleep(20 * time.Millisecond)

	hist.mu.Lock()
	count := len(hist.messages)
	hist.mu.Unlock()

	if count != 1 {
		t.Errorf("expected 1 message in history, got %d", count)
	}

	clientSide.Close()
	<-done
}

func TestHandleClient_RemovesClientOnDisconnect(t *testing.T) {
	serverSide, clientSide := net.Pipe()

	reg := NewRegistry()
	hist := &History{}

	c := &client.Client{Conn: serverSide}
	reg.Add(c)

	done := make(chan struct{})
	go func() {
		HandleClient(c, reg, hist)
		close(done)
	}()

	buf := make([]byte, 512)
	clientSide.Read(buf)
	fmt.Fprint(clientSide, "Theo\n")
	time.Sleep(20 * time.Millisecond)

	clientSide.Close()
	<-done

	if len(reg.All()) != 0 {
		t.Errorf("expected client to be removed from registry after disconnect")
	}
}

func TestHistory_EmptyReplay(t *testing.T) {
	h := &History{}

	serverSide, clientSide := net.Pipe()
	defer serverSide.Close()
	defer clientSide.Close()

	h.Replay(serverSide)

	clientSide.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	buf := make([]byte, 64)
	n, _ := clientSide.Read(buf)
	if n > 0 {
		t.Errorf("expected no data from empty history, got %q", string(buf[:n]))
	}
}
