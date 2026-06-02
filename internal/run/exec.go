package run

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// Exec is the production Runner backed by os/exec.
type Exec struct{}

// Run executes name with args, capturing stdout and stderr separately.
//
// A missing executable is reported as ErrNotFound (wrapped) so callers can
// degrade gracefully. Context cancellation/timeout terminates the process.
func (Exec) Run(ctx context.Context, name string, args ...string) (stdout, stderr []byte, err error) {
	if _, lookErr := exec.LookPath(name); lookErr != nil {
		return nil, nil, fmt.Errorf("%s: %w", name, ErrNotFound)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	return outBuf.Bytes(), errBuf.Bytes(), err
}
