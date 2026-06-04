package collect

import (
	"context"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// TCClass collects traffic-control shaping classes (htb/hfsc) per interface.
type TCClass struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewTCClass returns a TCClass collector using r and fs.
func NewTCClass(r run.Runner, fs sysfs.FS) *TCClass { return &TCClass{Runner: r, FS: fs} }

// Collect runs `tc -json class show dev <iface>` for each interface and keeps
// only rate-bearing (shaping) classes. Older iproute2 emits text rather than
// JSON for `class show`; that is reported once as a non-fatal warning.
func (c *TCClass) Collect(ctx context.Context) ([]model.TCClass, []model.Warning) {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil, nil
	}

	var out []model.TCClass
	var warnings []model.Warning
	textWarned := false
	for _, e := range entries {
		name := e.Name()
		stdout, _, rerr := c.Runner.Run(ctx, tcCmd, "-json", "class", "show", "dev", name)
		if rerr != nil {
			if run.IsNotFound(rerr) {
				return nil, []model.Warning{{Source: "tc", Message: "tc not found; class data unavailable"}}
			}
			continue
		}
		classes, perr := parse.TCClasses(stdout)
		if perr != nil {
			if !textWarned && strings.HasPrefix(strings.TrimSpace(string(stdout)), "class ") {
				warnings = append(warnings, model.Warning{Source: "tc", Message: "this iproute2 does not support JSON for `tc class`; class data unavailable"})
				textWarned = true
			}
			continue
		}
		for _, cl := range classes {
			if cl.Rate == 0 {
				continue // skip non-shaping classes (e.g. mq)
			}
			cl.Dev = name
			out = append(out, cl)
		}
	}
	return out, warnings
}
