package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestInterfacesCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: string(testutil.Fixture(t, "ip/address.json"))},
		"ip", "-detail", "-json", "address", "show",
	)

	ifaces, warnings := NewInterfaces(fake).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(ifaces) != 7 {
		t.Fatalf("interfaces = %d, want 7", len(ifaces))
	}
}

func TestInterfacesCollectMissingIP(t *testing.T) {
	ifaces, warnings := NewInterfaces(run.NewFake()).Collect(context.Background())
	if ifaces != nil {
		t.Errorf("interfaces = %v, want nil", ifaces)
	}
	if len(warnings) != 1 || !warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one fatal warning", warnings)
	}
	if warnings[0].Source != "ip" {
		t.Errorf("warning source = %q, want ip", warnings[0].Source)
	}
}
