package client

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestReadName_ValidOnFirstTry(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer serverSide.Close()
	defer clientSide.Close()

	go func() {
		fmt.Fprint(clientSide, "Yenlik\n")
	}()

	name, err := ReadName(serverSide)
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

	go func() {
		fmt.Fprint(clientSide, "\n")
		buf := make([]byte, len("[ENTER YOUR NAME]: "))
		clientSide.Read(buf)
		fmt.Fprint(clientSide, "Lee\n")
	}()

	name, err := ReadName(serverSide)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "Lee" {
		t.Errorf("expected 'Lee', got %q", name)
	}
}

func TestNotifyJoin_SendsToOthers(t *testing.T) {
	existing1Server, existing1Client := net.Pipe()
	existing2Server, existing2Client := net.Pipe()
	joinerServer, joinerClient := net.Pipe()
	defer existing1Server.Close()
	defer existing1Client.Close()
	defer existing2Server.Close()
	defer existing2Client.Close()
	defer joinerServer.Close()
	defer joinerClient.Close()

	existing1 := &Client{Conn: existing1Server, Name: "Yenlik"}
	existing2 := &Client{Conn: existing2Server, Name: "Lee"}
	joiner := &Client{Conn: joinerServer, Name: "Bob"}
	all := []*Client{existing1, existing2, joiner}

	done := make(chan string, 2)
	go func() {
		buf := make([]byte, 64)
		n, _ := existing1Client.Read(buf)
		done <- string(buf[:n])
	}()
	go func() {
		buf := make([]byte, 64)
		n, _ := existing2Client.Read(buf)
		done <- string(buf[:n])
	}()

	NotifyJoin(all, joiner)

	for range 2 {
		if got := <-done; got != "Bob has joined our chat...\n" {
			t.Errorf("expected join message, got %q", got)
		}
	}

	// Joiner should NOT receive the notification
	joinerClient.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	buf := make([]byte, 64)
	if n, _ := joinerClient.Read(buf); n > 0 {
		t.Errorf("joiner received own join message: %q", string(buf[:n]))
	}
}

func TestNotifyLeave_SendsToOthers(t *testing.T) {
	remaining1Server, remaining1Client := net.Pipe()
	remaining2Server, remaining2Client := net.Pipe()
	leaverServer, leaverClient := net.Pipe()
	defer remaining1Server.Close()
	defer remaining1Client.Close()
	defer remaining2Server.Close()
	defer remaining2Client.Close()
	defer leaverServer.Close()
	defer leaverClient.Close()

	remaining1 := &Client{Conn: remaining1Server, Name: "Yenlik"}
	remaining2 := &Client{Conn: remaining2Server, Name: "Lee"}
	leaver := &Client{Conn: leaverServer, Name: "Bob"}
	all := []*Client{remaining1, remaining2, leaver}

	done := make(chan string, 2)
	go func() {
		buf := make([]byte, 64)
		n, _ := remaining1Client.Read(buf)
		done <- string(buf[:n])
	}()
	go func() {
		buf := make([]byte, 64)
		n, _ := remaining2Client.Read(buf)
		done <- string(buf[:n])
	}()

	NotifyLeave(all, leaver)

	for range 2 {
		if got := <-done; got != "Bob has left our chat...\n" {
			t.Errorf("expected leave message, got %q", got)
		}
	}

	// Leaver should NOT receive the notification
	leaverClient.SetReadDeadline(time.Now().Add(20 * time.Millisecond))
	buf := make([]byte, 64)
	if n, _ := leaverClient.Read(buf); n > 0 {
		t.Errorf("leaver received own leave message: %q", string(buf[:n]))
	}
}

func TestSendWelcome_ContainsBanner(t *testing.T) {
	serverSide, clientSide := net.Pipe()
	defer serverSide.Close()
	defer clientSide.Close()

	go func() {
		SendWelcome(serverSide)
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
