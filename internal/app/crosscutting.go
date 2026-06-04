package app

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// netnsReexecEnv marks a process that has already re-execed into a namespace,
// preventing infinite recursion.
const netnsReexecEnv = "IFSCOPE_NETNS_ACTIVE"

// withWatch runs once, or repeatedly on the --watch interval (clearing the
// screen between runs) until interrupted. In watch mode, per-run errors are
// ignored so transient failures don't stop the loop.
func (o *Options) withWatch(run func() error) error {
	if o.Watch <= 0 {
		return run()
	}
	for {
		fmt.Print("\033[H\033[2J") // cursor home + clear screen
		_ = run()
		time.Sleep(o.Watch)
	}
}

// reexecNetns, when --netns is set and not already active, re-executes ifscope
// inside the target namespace via `ip netns exec`, forwarding stdio and exit
// code. Entering a namespace requires root, so this typically needs sudo.
func (o *Options) reexecNetns() error {
	if o.Netns == "" || os.Getenv(netnsReexecEnv) == "1" {
		return nil
	}

	exe, err := os.Executable()
	if err != nil {
		return &ExitCodeError{Code: ExitError, Err: fmt.Errorf("resolve executable: %w", err)}
	}

	args := append([]string{"netns", "exec", o.Netns, exe}, os.Args[1:]...)
	cmd := exec.Command("ip", args...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	cmd.Env = append(os.Environ(), netnsReexecEnv+"=1")

	err = cmd.Run()
	if err == nil {
		os.Exit(ExitOK)
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		os.Exit(exitErr.ExitCode())
	}
	// ip missing or could not enter the namespace.
	fmt.Fprintf(os.Stderr, "ifscope: cannot enter netns %q: %v\n", o.Netns, err)
	os.Exit(ExitPermDenied)
	return nil
}
