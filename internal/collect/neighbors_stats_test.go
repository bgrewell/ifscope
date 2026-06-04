package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
)

func TestNeighborsCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: `[{"dst":"10.0.0.1","dev":"eth0","lladdr":"aa:bb:cc:dd:ee:ff","state":["REACHABLE"]}]`},
		"ip", "-json", "neigh", "show",
	)
	ns, warnings := NewNeighbors(fake).Collect(context.Background())
	if len(warnings) != 0 || len(ns) != 1 || ns[0].Dst != "10.0.0.1" {
		t.Fatalf("neighbors = %v, warnings = %v", ns, warnings)
	}
}

func TestNeighborsCollectMissingIP(t *testing.T) {
	_, warnings := NewNeighbors(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || !warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one fatal", warnings)
	}
}

func TestStatsCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: `[{"ifname":"eth0","stats64":{"rx":{"bytes":100},"tx":{"bytes":200}}}]`},
		"ip", "-s", "-j", "link", "show",
	)
	stats, warnings := NewStats(fake).Collect(context.Background())
	if len(warnings) != 0 || len(stats) != 1 || stats[0].RxBytes != 100 || stats[0].TxBytes != 200 {
		t.Fatalf("stats = %v, warnings = %v", stats, warnings)
	}
}

func TestStatsCollectMissingIP(t *testing.T) {
	_, warnings := NewStats(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || !warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one fatal", warnings)
	}
}
