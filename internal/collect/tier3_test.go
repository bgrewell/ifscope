package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

func TestQdiscCollectRootsOnly(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: `[
		  {"kind":"mq","handle":"0:","dev":"eth0","root":true},
		  {"kind":"pfifo_fast","handle":"0:","dev":"eth0","parent":":1"},
		  {"kind":"fq_codel","handle":"8001:","dev":"eth1","root":true}
		]`},
		"tc", "-json", "qdisc", "show",
	)
	qs, warnings := NewQdisc(fake).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(qs) != 2 {
		t.Fatalf("qdiscs = %d, want 2 (roots only)", len(qs))
	}
	for _, q := range qs {
		if !q.Root {
			t.Errorf("non-root qdisc leaked: %+v", q)
		}
	}
}

func TestQdiscCollectMissingTC(t *testing.T) {
	_, warnings := NewQdisc(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || warnings[0].Source != "tc" || warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one non-fatal tc", warnings)
	}
}

func TestOffloadsCollect(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{"eth0", "lo"}
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: "Features for eth0:\nrx-checksumming: on\ntcp-segmentation-offload: off [fixed]\n"},
		"ethtool", "-k", "eth0",
	)
	// lo is skipped, so no ethtool call is registered for it.

	off, warnings := NewOffloads(fake, fs).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(off) != 1 || off[0].Name != "eth0" {
		t.Fatalf("offloads = %+v, want one for eth0 (lo skipped)", off)
	}
	if off[0].Features["rx-checksumming"] != "on" || off[0].Features["tcp-segmentation-offload"] != "off [fixed]" {
		t.Errorf("features = %v", off[0].Features)
	}
}

func TestOffloadsCollectMissingEthtool(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{"eth0"}
	_, warnings := NewOffloads(run.NewFake(), fs).Collect(context.Background())
	if len(warnings) != 1 || warnings[0].Source != "ethtool" {
		t.Fatalf("warnings = %v, want one ethtool warning", warnings)
	}
}
