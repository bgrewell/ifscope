package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestEthtoolEnrich(t *testing.T) {
	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "ethtool/ice-i.txt"))}, "ethtool", "-i", "eth0")
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "ethtool/ice-settings.txt"))}, "ethtool", "eth0")

	ifaces := []model.Interface{{Name: "eth0", Type: model.TypePhysical}}
	if w := NewEthtool(fake).Enrich(context.Background(), ifaces); w != nil {
		t.Fatalf("unexpected warnings: %v", w)
	}

	got := ifaces[0]
	if got.Driver != "ice" {
		t.Errorf("driver = %q, want ice", got.Driver)
	}
	if got.Speed != "100 Gb/s" {
		t.Errorf("speed = %q, want 100 Gb/s", got.Speed)
	}
	if got.Port != "Direct Attach Copper" {
		t.Errorf("port = %q", got.Port)
	}
	if got.Bus != "0000:17:00.0" {
		t.Errorf("bus = %q, want 0000:17:00.0", got.Bus)
	}
}

func TestEthtoolEnrichMissingBinary(t *testing.T) {
	ifaces := []model.Interface{{Name: "eth0"}}
	w := NewEthtool(run.NewFake()).Enrich(context.Background(), ifaces)
	if len(w) != 1 || w[0].Source != "ethtool" || w[0].Fatal {
		t.Fatalf("warnings = %v, want one non-fatal ethtool warning", w)
	}
}

func TestIsPCIBus(t *testing.T) {
	cases := map[string]bool{
		"0000:17:00.0": true,
		"0000:25:00.3": true,
		"tap":          false,
		"":             false,
		"17:00.0":      false,
	}
	for in, want := range cases {
		if got := isPCIBus(in); got != want {
			t.Errorf("isPCIBus(%q) = %v, want %v", in, got, want)
		}
	}
}
