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

func TestIPRoutesMultipath(t *testing.T) {
	// An ECMP route carries next-hops with no top-level gateway/dev.
	data := []byte(`[{
	  "dst":"default","protocol":"static",
	  "nexthops":[
	    {"gateway":"10.0.0.1","dev":"eth0","weight":1},
	    {"gateway":"10.0.0.2","dev":"eth1","weight":2}
	  ]
	}]`)
	routes, err := IPRoutes(data)
	if err != nil {
		t.Fatal(err)
	}
	r := routes[0]
	if r.Gateway != "" || r.Dev != "" {
		t.Errorf("multipath route should have empty top-level gateway/dev, got %q/%q", r.Gateway, r.Dev)
	}
	if len(r.NextHops) != 2 {
		t.Fatalf("nexthops = %d, want 2", len(r.NextHops))
	}
	if r.NextHops[0].Gateway != "10.0.0.1" || r.NextHops[0].Dev != "eth0" || r.NextHops[0].Weight != 1 {
		t.Errorf("nexthop0 = %+v", r.NextHops[0])
	}
	if r.Family != "inet" {
		t.Errorf("family = %q", r.Family)
	}
}

func TestIPRoutesMultipathIPv6Family(t *testing.T) {
	data := []byte(`[{"dst":"default","nexthops":[{"gateway":"fe80::1","dev":"eth0"}]}]`)
	routes, _ := IPRoutes(data)
	if routes[0].Family != "inet6" {
		t.Errorf("family = %q, want inet6 (from next-hop)", routes[0].Family)
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
