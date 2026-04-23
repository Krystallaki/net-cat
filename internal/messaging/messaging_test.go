package messaging

import (
	"net"
	"testing"
	"time"
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
	defer serverSide.Close()
	defer clientSide.Close()

	done := make(chan string, 1)
	go func() {
		buf := make([]byte, 256)
		n, _ := clientSide.Read(buf)
		done <- string(buf[:n])
	}()

	h.Replay(serverSide)
	serverSide.Close()

	got := <-done
	want := "[2020-01-20 15:48:41][Yenlik]:hello\n[2020-01-20 15:48:42][Yenlik]:world\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
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
