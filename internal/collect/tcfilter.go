package collect

import (
	"context"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// TCFilter collects traffic-control filters (classifiers) per interface.
type TCFilter struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewTCFilter returns a TCFilter collector using r and fs.
func NewTCFilter(r run.Runner, fs sysfs.FS) *TCFilter { return &TCFilter{Runner: r, FS: fs} }

// Collect runs `tc -json filter show dev <iface>` for each interface. Older
// iproute2 emits text rather than JSON for `filter show`; that is reported once
// as a non-fatal warning.
func (c *TCFilter) Collect(ctx context.Context) ([]model.TCFilter, []model.Warning) {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil, nil
	}

	var out []model.TCFilter
	var warnings []model.Warning
	textWarned := false
	for _, e := range entries {
		name := e.Name()
		stdout, _, rerr := c.Runner.Run(ctx, tcCmd, "-json", "filter", "show", "dev", name)
		if rerr != nil {
			if run.IsNotFound(rerr) {
				return nil, []model.Warning{{Source: "tc", Message: "tc not found; filter data unavailable"}}
			}
			continue
		}
		filters, perr := parse.TCFilters(stdout)
		if perr != nil {
			if !textWarned && strings.HasPrefix(strings.TrimSpace(string(stdout)), "filter ") {
				warnings = append(warnings, model.Warning{Source: "tc", Message: "this iproute2 does not support JSON for `tc filter`; filter data unavailable"})
				textWarned = true
			}
			continue
		}
		for _, f := range filters {
			f.Dev = name
			out = append(out, f)
		}
	}
	return out, warnings
}
