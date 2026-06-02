package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestIPRoutes(t *testing.T) {
	routes, err := IPRoutes(testutil.Fixture(t, "ip/route.json"))
	if err != nil {
		t.Fatalf("IPRoutes: %v", err)
	}
	if len(routes) == 0 {
		t.Fatal("no routes parsed")
	}

	def := routes[0]
	if def.Dst != "default" {
		t.Errorf("first route dst = %q, want default", def.Dst)
	}
	if def.Gateway == "" {
		t.Errorf("default route missing gateway")
	}
	if def.Family != "inet" {
		t.Errorf("family = %q, want inet", def.Family)
	}

	// A kernel route should carry its preferred source as Src.
	var foundSrc bool
	for _, r := range routes {
		if r.Protocol == "kernel" && r.Src != "" {
			foundSrc = true
		}
	}
	if !foundSrc {
		t.Errorf("expected a kernel route with a source address")
	}
}

func TestIPRoutesFamilyInference(t *testing.T) {
	data := []byte(`[{"dst":"::/0","gateway":"fe80::1","dev":"eth0","protocol":"ra"}]`)
	routes, err := IPRoutes(data)
	if err != nil {
		t.Fatal(err)
	}
	if routes[0].Family != "inet6" {
		t.Errorf("family = %q, want inet6", routes[0].Family)
	}
}
