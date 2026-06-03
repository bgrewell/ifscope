package collect

import (
	"context"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Rules collects routing policy rules (`ip rule`).
type Rules struct {
	Runner run.Runner
}

// NewRules returns a Rules collector using r.
func NewRules(r run.Runner) *Rules { return &Rules{Runner: r} }

// Collect runs `ip -json rule` and parses the policy rules.
func (c *Rules) Collect(ctx context.Context) ([]model.Rule, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-json", "rule")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ip", Message: "ip command not found; policy rules unavailable", Fatal: true}}
		}
		return nil, []model.Warning{{
			Source:  "ip",
			Message: fmt.Sprintf("ip rule failed: %v: %s", err, strings.TrimSpace(string(stderr))),
		}}
	}
	rules, perr := parse.IPRules(stdout)
	if perr != nil {
		return nil, []model.Warning{{Source: "ip", Message: perr.Error()}}
	}
	return rules, nil
}
