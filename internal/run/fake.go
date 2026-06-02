package run

import (
	"context"
	"fmt"
	"strings"
)

// Fake is a test Runner that returns canned results keyed by the full command
// line ("name arg1 arg2 ..."). Unkeyed commands return ErrNotFound, modeling a
// host where the tool is absent.
type Fake struct {
	// Responses maps a command line to its canned result.
	Responses map[string]FakeResult
	// Calls records every command line passed to Run, in order.
	Calls []string
}

// FakeResult is a canned command outcome.
type FakeResult struct {
	Stdout string
	Stderr string
	Err    error
}

// NewFake returns an empty Fake ready for Set.
func NewFake() *Fake {
	return &Fake{Responses: map[string]FakeResult{}}
}

// Set registers a canned result for the given command line.
func (f *Fake) Set(stdout string, name string, args ...string) *Fake {
	f.Responses[key(name, args)] = FakeResult{Stdout: stdout}
	return f
}

// SetResult registers a full canned result (stdout/stderr/err).
func (f *Fake) SetResult(r FakeResult, name string, args ...string) *Fake {
	f.Responses[key(name, args)] = r
	return f
}

// Run returns the canned result for the command line, recording the call.
func (f *Fake) Run(_ context.Context, name string, args ...string) ([]byte, []byte, error) {
	k := key(name, args)
	f.Calls = append(f.Calls, k)
	r, ok := f.Responses[k]
	if !ok {
		return nil, nil, fmt.Errorf("%s: %w", name, ErrNotFound)
	}
	return []byte(r.Stdout), []byte(r.Stderr), r.Err
}

func key(name string, args []string) string {
	if len(args) == 0 {
		return name
	}
	return name + " " + strings.Join(args, " ")
}
