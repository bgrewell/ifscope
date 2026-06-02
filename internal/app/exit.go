package app

import "errors"

// errTestFailed marks a connectivity-test failure for exit-code mapping.
var errTestFailed = errors.New("one or more connectivity tests failed")

// Exit codes per the REL-1 spec. Inspection commands return 0 even when
// optional data is missing; only explicit failures use non-zero codes.
const (
	ExitOK         = 0
	ExitError      = 1
	ExitUsage      = 2
	ExitTestFailed = 10
	ExitDepMissing = 20
	ExitPermDenied = 30
)

// ExitError carries an exit code alongside an error message so the top-level
// Execute can map command failures to process exit codes.
type ExitCodeError struct {
	Code int
	Err  error
}

func (e *ExitCodeError) Error() string { return e.Err.Error() }
func (e *ExitCodeError) Unwrap() error { return e.Err }
