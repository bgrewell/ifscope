package collect

import (
	"context"
	"encoding/json"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Netns collects named network namespaces via iproute2.
type Netns struct {
	Runner run.Runner
}

// NewNetns returns a Netns collector using r.
func NewNetns(r run.Runner) *Netns { return &Netns{Runner: r} }

// Collect lists network namespaces and, where permitted, counts the interfaces
// inside each. A missing ip is a non-fatal warning (namespaces are optional
// context); per-namespace counts that require privileges are left unknown.
func (c *Netns) Collect(ctx context.Context) ([]model.Netns, []model.Warning) {
	stdout, _, err := c.Runner.Run(ctx, ipCmd, "netns", "list")
	if err != nil {
		if run.IsNotFound(err) {
			return nil, []model.Warning{{Source: "ip", Message: "ip command not found; namespace data unavailable"}}
		}
		// `ip netns list` can fail if the netns directory is absent; treat as
		// "no namespaces" rather than an error.
		return nil, nil
	}

	namespaces := parse.NetnsList(stdout)
	for i := range namespaces {
		if n := c.interfaceCount(ctx, namespaces[i].Name); n >= 0 {
			count := n
			namespaces[i].Interfaces = &count
		}
	}
	return namespaces, nil
}

// interfaceCount returns the number of links inside a namespace, or -1 when it
// cannot be determined (e.g. without privileges to enter the namespace).
func (c *Netns) interfaceCount(ctx context.Context, name string) int {
	out, _, err := c.Runner.Run(ctx, ipCmd, "-n", name, "-json", "link", "show")
	if err != nil {
		return -1
	}
	var links []struct{}
	if json.Unmarshal(out, &links) != nil {
		return -1
	}
	return len(links)
}
