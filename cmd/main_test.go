package main

import "testing"

// TestDefaultPort_NoArgs verifies that the server defaults to port 8989 when no argument is given.
func TestDefaultPort_NoArgs(t *testing.T) {
	if got := resolvePort([]string{}); got != "8989" {
		t.Errorf("expected default port 8989, got %q", got)
	}
}

// TestDefaultPort_WithArg verifies that the provided port is used when an argument is given.
func TestDefaultPort_WithArg(t *testing.T) {
	if got := resolvePort([]string{"9000"}); got != "9000" {
		t.Errorf("expected port 9000, got %q", got)
	}
}

// TestParseArgs_TooManyArgs verifies that more than one argument is rejected.
func TestParseArgs_TooManyArgs(t *testing.T) {
	_, ok := parseArgs([]string{"9000", "9001"})
	if ok {
		t.Error("expected parseArgs to return false for more than one argument")
	}
}
