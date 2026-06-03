package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestRoutesCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: string(testutil.Fixture(t, "ip/route.json"))},
		"ip", "-detail", "-json", "route", "show", "table", "all",
	)
	routes, warnings := NewRoutes(fake).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(routes) == 0 || routes[0].Dst != "default" {
		t.Fatalf("routes = %v", routes)
	}
	// Main-table routes (table omitted by iproute2) normalize to "main".
	if routes[0].Table != "main" {
		t.Errorf("table = %q, want main", routes[0].Table)
	}
}

func TestRoutesCollectMissingIP(t *testing.T) {
	_, warnings := NewRoutes(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || !warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one fatal", warnings)
	}
}

func TestDNSCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: string(testutil.Fixture(t, "resolvectl/status.txt"))},
		"resolvectl", "status",
	)
	dns, warnings := NewDNS(fake).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(dns) != 3 {
		t.Fatalf("dns blocks = %d, want 3", len(dns))
	}
}

func TestDNSCollectMissingResolvectl(t *testing.T) {
	dns, warnings := NewDNS(run.NewFake()).Collect(context.Background())
	if dns != nil {
		t.Errorf("dns = %v, want nil", dns)
	}
	if len(warnings) != 1 || warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one non-fatal", warnings)
	}
}
