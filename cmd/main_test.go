package main

import "testing"

func TestDefaultPort_NoArgs(t *testing.T) {
	if got := resolvePort([]string{}); got != "8989" {
		t.Errorf("expected default port 8989, got %q", got)
	}
}

func TestDefaultPort_WithArg(t *testing.T) {
	if got := resolvePort([]string{"9000"}); got != "9000" {
		t.Errorf("expected port 9000, got %q", got)
	}
}

func TestParseArgs_TooManyArgs(t *testing.T) {
	_, ok := parseArgs([]string{"9000", "9001"})
	if ok {
		t.Error("expected parseArgs to return false for more than one argument")
	}
}
