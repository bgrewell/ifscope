package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
)

func TestDevlinkCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: `{"port":{"pci/0000:25:00.0/0":{"type":"eth","netdev":"eth0","flavour":"physical","lanes":4}}}`},
		"devlink", "-j", "port", "show",
	)
	ports, warnings := NewDevlink(fake).Collect(context.Background())
	if len(warnings) != 0 || len(ports) != 1 || ports[0].Netdev != "eth0" {
		t.Fatalf("ports=%v warnings=%v", ports, warnings)
	}
}

func TestDevlinkCollectMissing(t *testing.T) {
	_, warnings := NewDevlink(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || warnings[0].Source != "devlink" || warnings[0].Fatal {
		t.Fatalf("warnings=%v, want one non-fatal devlink", warnings)
	}
}

func TestFDBCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: `[{"mac":"aa:bb:cc:dd:ee:ff","ifname":"eth0","master":"br0","state":"permanent"}]`},
		"bridge", "-json", "fdb", "show",
	)
	entries, warnings := NewFDB(fake).Collect(context.Background())
	if len(warnings) != 0 || len(entries) != 1 || entries[0].Master != "br0" {
		t.Fatalf("entries=%v warnings=%v", entries, warnings)
	}
}

func TestFDBCollectMissing(t *testing.T) {
	_, warnings := NewFDB(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || warnings[0].Source != "bridge" || warnings[0].Fatal {
		t.Fatalf("warnings=%v, want one non-fatal bridge", warnings)
	}
}

func TestSocketsCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: "tcp LISTEN 0 4096 *:22 *:*\n"},
		"ss", "-tulpnH",
	)
	socks, warnings := NewSockets(fake).Collect(context.Background())
	if len(warnings) != 0 || len(socks) != 1 || socks[0].LocalPort != "22" {
		t.Fatalf("socks=%v warnings=%v", socks, warnings)
	}
}

func TestSocketsCollectMissing(t *testing.T) {
	_, warnings := NewSockets(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || warnings[0].Source != "ss" || warnings[0].Fatal {
		t.Fatalf("warnings=%v, want one non-fatal ss", warnings)
	}
}
