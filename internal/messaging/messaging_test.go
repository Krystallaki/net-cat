package messaging

import (
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
