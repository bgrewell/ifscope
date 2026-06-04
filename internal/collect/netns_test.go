package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
)

func TestNetnsCollect(t *testing.T) {
	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: "myns (id: 0)\nother\n"}, "ip", "netns", "list")
	fake.SetResult(run.FakeResult{Stdout: `[{},{}]`}, "ip", "-n", "myns", "-json", "link", "show")
	// "other" can't be entered (permission) -> count unknown.
	fake.SetResult(run.FakeResult{Err: errExit}, "ip", "-n", "other", "-json", "link", "show")

	ns, warnings := NewNetns(fake).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(ns) != 2 {
		t.Fatalf("namespaces = %d, want 2", len(ns))
	}
	if ns[0].Interfaces == nil || *ns[0].Interfaces != 2 {
		t.Errorf("myns interfaces = %v, want 2", ns[0].Interfaces)
	}
	if ns[1].Interfaces != nil {
		t.Errorf("other interfaces = %v, want nil (unknown)", ns[1].Interfaces)
	}
}

func TestNetnsCollectMissingIP(t *testing.T) {
	ns, warnings := NewNetns(run.NewFake()).Collect(context.Background())
	if ns != nil {
		t.Errorf("ns = %v, want nil", ns)
	}
	if len(warnings) != 1 || warnings[0].Source != "ip" || warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one non-fatal ip warning", warnings)
	}
}
