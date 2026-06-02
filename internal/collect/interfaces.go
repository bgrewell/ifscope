// Package collect contains collectors that orchestrate the command runner and
// sysfs reads, feed raw output to parsers, and return typed data plus warnings.
// Collectors never print and never abort the program; missing optional tooling
// is reported as warnings so the report degrades gracefully.
package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const ipCmd = "ip"

// Interfaces collects interface and VLAN data via iproute2.
type Interfaces struct {
	Runner run.Runner
}

// NewInterfaces returns an Interfaces collector using r.
func NewInterfaces(r run.Runner) *Interfaces { return &Interfaces{Runner: r} }

// Collect runs `ip -detail -json address show` and parses it. The returned
// slice includes VLAN interfaces; callers partition by Type. A missing or
// failing `ip` is fatal for this collector and reported as such.
func (c *Interfaces) Collect(ctx context.Context) ([]model.Interface, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-detail", "-json", "address", "show")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{
				Source:  "ip",
				Message: "ip command not found; interface data unavailable",
				Fatal:   true,
			}}
		}
		return nil, []model.Warning{{
			Source:  "ip",
			Message: fmt.Sprintf("ip address show failed: %v: %s", err, strings.TrimSpace(string(stderr))),
			Fatal:   true,
		}}
	}

	ifaces, perr := parse.IPAddresses(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "ip", Message: perr.Error(), Fatal: true}}
	}
	return ifaces, nil
}
