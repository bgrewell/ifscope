package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
)

func TestRoutesRenderMultipath(t *testing.T) {
	routes := []model.Route{
		{Dst: "10.0.0.0/24", Gateway: "192.168.1.1", Dev: "eth0", Protocol: "static"},
		{
			Dst:      "default",
			Protocol: "static",
			NextHops: []model.RouteNextHop{
				{Gateway: "10.0.0.1", Dev: "eth0", Weight: 1},
				{Gateway: "10.0.0.2", Dev: "eth1", Weight: 2},
			},
		},
	}

	var buf bytes.Buffer
	Options{Color: NewColor("never", false, nil)}.Routes(&buf, routes)
	out := buf.String()

	// Single-path route renders normally.
	if !strings.Contains(out, "192.168.1.1") {
		t.Errorf("missing single-path gateway in:\n%s", out)
	}
	// Multipath route shows both next-hops with their devices and weights.
	for _, want := range []string{"10.0.0.1 (w1)", "10.0.0.2 (w2)", "eth0", "eth1"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%s", want, out)
		}
	}
}
