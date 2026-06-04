package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

func TestQueuesCollect(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{"eth0", "lo"}
	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: "Channel parameters for eth0:\nPre-set maximums:\nCombined:\t192\nCurrent hardware settings:\nCombined:\t48\n"}, "ethtool", "-l", "eth0")
	fake.SetResult(run.FakeResult{Stdout: "Ring parameters for eth0:\nPre-set maximums:\nRX:\t8160\nTX:\t8160\nCurrent hardware settings:\nRX:\t2048\nTX:\t256\n"}, "ethtool", "-g", "eth0")

	qs, warnings := NewQueues(fake, fs).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(qs) != 1 || qs[0].Name != "eth0" {
		t.Fatalf("queues = %+v (lo should be skipped)", qs)
	}
	if qs[0].Combined.Current != 48 || qs[0].Combined.Max != 192 {
		t.Errorf("combined = %+v", qs[0].Combined)
	}
	if qs[0].RxRing.Current != 2048 || qs[0].TxRing.Current != 256 {
		t.Errorf("rings = rx %+v tx %+v", qs[0].RxRing, qs[0].TxRing)
	}
}

func TestIRQCollect(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{"eth0", "lo"}
	fs.Files["/proc/interrupts"] = " 104:  0  0  IR-PCI-MSIX  2-edge  ice-eth0-TxRx-0\n NMI:  1  2  Non-maskable interrupts\n"
	fs.Files["/proc/irq/104/smp_affinity_list"] = "0-3\n"

	irqs, _ := NewIRQ(fs).Collect()
	if len(irqs) != 1 {
		t.Fatalf("irqs = %d, want 1 (NMI skipped, matched to eth0)", len(irqs))
	}
	if irqs[0].Device != "eth0" || irqs[0].IRQ != 104 || irqs[0].CPUs != "0-3" {
		t.Errorf("irq = %+v", irqs[0])
	}
}

func TestMulticastCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: `[{"ifname":"eth0","maddr":[{"family":"inet","address":"224.0.0.1"}]}]`},
		"ip", "-json", "maddr", "show",
	)
	groups, warnings := NewMulticast(fake).Collect(context.Background())
	if len(warnings) != 0 || len(groups) != 1 || groups[0].Address != "224.0.0.1" {
		t.Fatalf("groups=%v warnings=%v", groups, warnings)
	}
}
