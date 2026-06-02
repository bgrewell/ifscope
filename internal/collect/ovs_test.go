package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/testutil"
)

// ovsFixtures registers the three ovs-vsctl table responses on a fake under the
// given command name ("ovs-vsctl" direct, or with a sudo prefix).
func setOVS(fake *run.Fake, res func(table, cols string) run.FakeResult) {
	fake.SetResult(res("Bridge", parse.OVSBridgeColumns), "ovs-vsctl", "--format=json", "--columns="+parse.OVSBridgeColumns, "list", "Bridge")
	fake.SetResult(res("Port", parse.OVSPortColumns), "ovs-vsctl", "--format=json", "--columns="+parse.OVSPortColumns, "list", "Port")
	fake.SetResult(res("Interface", parse.OVSInterfaceColumns), "ovs-vsctl", "--format=json", "--columns="+parse.OVSInterfaceColumns, "list", "Interface")
}

func fixtureFor(t *testing.T, table string) string {
	switch table {
	case "Bridge":
		return string(testutil.Fixture(t, "ovs/bridge.json"))
	case "Port":
		return string(testutil.Fixture(t, "ovs/port.json"))
	default:
		return string(testutil.Fixture(t, "ovs/interface.json"))
	}
}

func TestOVSCollectUnprivileged(t *testing.T) {
	fake := run.NewFake()
	setOVS(fake, func(table, _ string) run.FakeResult {
		return run.FakeResult{Stdout: fixtureFor(t, table)}
	})

	ovs, warnings := NewOVS(fake, false).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if ovs == nil || len(ovs.Bridges) != 1 {
		t.Fatalf("ovs = %+v", ovs)
	}
}

func TestOVSCollectAutoSudo(t *testing.T) {
	fake := run.NewFake()
	// Direct ovs-vsctl calls fail (permission); sudo variants succeed.
	for _, tc := range []struct{ table, cols string }{
		{"Bridge", parse.OVSBridgeColumns},
		{"Port", parse.OVSPortColumns},
		{"Interface", parse.OVSInterfaceColumns},
	} {
		fake.SetResult(run.FakeResult{Stderr: "database connection failed (Permission denied)", Err: errExit},
			"ovs-vsctl", "--format=json", "--columns="+tc.cols, "list", tc.table)
		fake.SetResult(run.FakeResult{Stdout: fixtureFor(t, tc.table)},
			"sudo", "-n", "ovs-vsctl", "--format=json", "--columns="+tc.cols, "list", tc.table)
	}

	ovs, warnings := NewOVS(fake, false).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if ovs == nil || len(ovs.Bridges) != 1 {
		t.Fatalf("ovs = %+v", ovs)
	}
}

func TestOVSCollectNoSudoDenied(t *testing.T) {
	fake := run.NewFake()
	setOVS(fake, func(_, _ string) run.FakeResult {
		return run.FakeResult{Stderr: "Permission denied", Err: errExit}
	})

	ovs, warnings := NewOVS(fake, true).Collect(context.Background())
	if ovs != nil {
		t.Errorf("ovs = %+v, want nil", ovs)
	}
	if len(warnings) != 1 || warnings[0].Source != "ovs" {
		t.Fatalf("warnings = %v, want one ovs warning", warnings)
	}
}

func TestOVSCollectMissingBinary(t *testing.T) {
	ovs, warnings := NewOVS(run.NewFake(), false).Collect(context.Background())
	if ovs != nil {
		t.Errorf("ovs = %+v, want nil", ovs)
	}
	if len(warnings) != 1 || warnings[0].Source != "ovs-vsctl" {
		t.Fatalf("warnings = %v, want ovs-vsctl not-found warning", warnings)
	}
}
