package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestResolvectl(t *testing.T) {
	dns := Resolvectl(testutil.Fixture(t, "resolvectl/status.txt"))
	if len(dns) != 3 {
		t.Fatalf("blocks = %d, want 3", len(dns))
	}

	by := make(map[string]model.DNS, len(dns))
	for _, d := range dns {
		by[d.Link] = d
	}

	t.Run("global wrapped servers and domains", func(t *testing.T) {
		g := by["global"]
		want := []string{"1.1.1.1", "8.8.8.8", "2606:4700:4700::1111"}
		if len(g.Servers) != 3 {
			t.Fatalf("servers = %v, want %v", g.Servers, want)
		}
		for i, s := range want {
			if g.Servers[i] != s {
				t.Errorf("server[%d] = %q, want %q", i, g.Servers[i], s)
			}
		}
		if len(g.Domains) != 2 || g.Domains[0] != "example.com" {
			t.Errorf("domains = %v", g.Domains)
		}
	})

	t.Run("link with default route and protocols", func(t *testing.T) {
		e := by["eth0"]
		if e.CurrentServer != "172.26.1.1" {
			t.Errorf("current = %q", e.CurrentServer)
		}
		if len(e.Servers) != 2 {
			t.Errorf("servers = %v, want 2", e.Servers)
		}
		if e.DefaultRoute == nil || !*e.DefaultRoute {
			t.Errorf("default route = %v, want true", e.DefaultRoute)
		}
		if e.LLMNR != "yes" {
			t.Errorf("llmnr = %q, want yes", e.LLMNR)
		}
		if e.DNSSEC != "yes/supported" {
			t.Errorf("dnssec = %q, want yes/supported", e.DNSSEC)
		}
	})

	t.Run("link with no servers", func(t *testing.T) {
		e := by["eth1"]
		if len(e.Servers) != 0 {
			t.Errorf("servers = %v, want none", e.Servers)
		}
		if e.DefaultRoute == nil || *e.DefaultRoute {
			t.Errorf("default route = %v, want false", e.DefaultRoute)
		}
	})
}
