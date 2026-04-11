package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
)

// net.Pipe() δημιουργεί δύο συνδεδεμένα net.Conn.
// Ό,τι γράψουμε στο serverSide φαίνεται στο clientSide.

func TestReadName_ValidOnFirstTry(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer serverSide.Close()
	defer clientSide.Close()

	go func() {
		fmt.Fprint(clientSide, "Yenlik\n")
	}()

	name, err := readName(serverSide)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "Yenlik" {
		t.Errorf("expected 'Yenlik', got %q", name)
	}
}

func TestReadName_EmptyThenValid(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer serverSide.Close()
	defer clientSide.Close()

	// Στέλνουμε κενή γραμμή πρώτα, μετά έγκυρο όνομα
	go func() {
		fmt.Fprint(clientSide, "\n")
		// Διαβάζουμε το re-prompt που έστειλε ο server
		buf := make([]byte, len("[ENTER YOUR NAME]: "))
		clientSide.Read(buf)
		fmt.Fprint(clientSide, "Lee\n")
	}()

	name, err := readName(serverSide)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "Lee" {
		t.Errorf("expected 'Lee', got %q", name)
	}
}

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
