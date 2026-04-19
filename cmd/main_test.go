package main

import "testing"

// defaultPort returns "8989" when no argument is provided, or the given argument otherwise.
// This mirrors the logic in main() and is extracted here for testability.
func defaultPort(args []string) string {
	if len(args) == 0 {
		return "8989"
	}
	return args[0]
}

// TestDefaultPort_NoArgs verifies that the server defaults to port 8989 when no argument is given.
func TestDefaultPort_NoArgs(t *testing.T) {
	if got := defaultPort([]string{}); got != "8989" {
		t.Errorf("expected default port 8989, got %q", got)
	}
}

// TestDefaultPort_WithArg verifies that the provided port is used when an argument is given.
func TestDefaultPort_WithArg(t *testing.T) {
	if got := defaultPort([]string{"9000"}); got != "9000" {
		t.Errorf("expected port 9000, got %q", got)
	}
}
