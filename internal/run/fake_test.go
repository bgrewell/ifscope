package run

import (
	"context"
	"testing"
)

func TestFakeReturnsCannedOutput(t *testing.T) {
	f := NewFake().Set("hello\n", "echo", "hello")

	stdout, _, err := f.Run(context.Background(), "echo", "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := string(stdout); got != "hello\n" {
		t.Fatalf("stdout = %q, want %q", got, "hello\n")
	}
	if len(f.Calls) != 1 || f.Calls[0] != "echo hello" {
		t.Fatalf("Calls = %v, want [echo hello]", f.Calls)
	}
}

func TestFakeUnknownCommandIsNotFound(t *testing.T) {
	f := NewFake()

	_, _, err := f.Run(context.Background(), "missing")
	if !IsNotFound(err) {
		t.Fatalf("err = %v, want not-found", err)
	}
}
