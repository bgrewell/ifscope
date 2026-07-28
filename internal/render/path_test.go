package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
)

func TestPathRendersDecision(t *testing.T) {
	path := model.Path{Candidates: []model.PathCandidate{{
		Address: "203.0.113.10",
		Family:  "inet",
		Route:   &model.Route{Table: "100", Src: "192.0.2.10", Gateway: "192.0.2.1", Dev: "vlan10"},
		Rule:    &model.Rule{Priority: 100},
		Neighbor: &model.Neighbor{
			State: "REACHABLE", LLAddr: "00:11:22:33:44:55",
		},
		Topology: []model.PathHop{
			{Name: "vlan10", Type: model.TypeVLAN},
			{Name: "eth0", Type: model.TypePhysical},
		},
		MTU: 1500,
	}}}
	var buf bytes.Buffer
	Options{Barebones: true}.Path(&buf, path)
	out := buf.String()
	for _, want := range []string{"203.0.113.10", "100/100", "192.0.2.1", "REACHABLE", "vlan10 (vlan) → eth0 (physical)", "1500"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}
