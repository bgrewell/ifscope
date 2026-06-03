package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/bgrewell/ifscope/internal/testutil"
)

// pciFS builds a fake PCI tree with four Ethernet devices: a kernel-bound NIC,
// a vfio-pci (DPDK) NIC with no netdev, an unbound NIC, and a non-Ethernet
// device that must be ignored.
func pciFS() *sysfs.Fake {
	fs := sysfs.NewFake()
	fs.Dirs[pciDevices] = []string{"0000:17:00.0", "0000:43:00.2", "0000:43:00.3", "0000:00:1f.0"}

	mk := func(bus, class, driver string) string {
		base := pciDevices + "/" + bus
		fs.Files[base+"/class"] = class + "\n"
		fs.Files[base+"/vendor"] = "0x8086\n"
		fs.Files[base+"/device"] = "0x1592\n"
		fs.Files[base+"/numa_node"] = "0\n"
		if driver != "" {
			fs.Links[base+"/driver"] = "../../../bus/pci/drivers/" + driver
		}
		return base
	}

	// kernel NIC with a netdev
	base := mk("0000:17:00.0", "0x020000", "ice")
	fs.Dirs[base+"/net"] = []string{"enp23s0np0"}
	// DPDK NIC bound to vfio-pci, no net dir
	mk("0000:43:00.2", "0x020000", "vfio-pci")
	// unbound NIC, no driver, no net
	mk("0000:43:00.3", "0x020000", "")
	// non-Ethernet device (ISA bridge) — ignored
	mk("0000:00:1f.0", "0x060100", "lpc_ich")

	return fs
}

func TestPCIeScanClassifiesBind(t *testing.T) {
	fake := run.NewFake() // lspci absent -> non-fatal warning
	ifaces := []model.Interface{{Name: "enp23s0np0", Bus: "0000:17:00.0", Driver: "ice"}}

	devices, warnings := NewPCIe(fake, pciFS()).Collect(context.Background(), ifaces)

	byBus := map[string]model.PCIDevice{}
	for _, d := range devices {
		byBus[d.Bus] = d
	}

	if len(devices) != 3 {
		t.Fatalf("devices = %d, want 3 (non-Ethernet excluded)", len(devices))
	}

	kernel := byBus["0000:17:00.0"]
	if kernel.Bind != "kernel" || kernel.Interface != "enp23s0np0" {
		t.Errorf("kernel device = %+v", kernel)
	}
	if kernel.NUMANode == nil || *kernel.NUMANode != 0 {
		t.Errorf("kernel numa = %v", kernel.NUMANode)
	}

	dpdk := byBus["0000:43:00.2"]
	if dpdk.Bind != "dpdk" || dpdk.Driver != "vfio-pci" || dpdk.Interface != "" {
		t.Errorf("dpdk device = %+v, want bind=dpdk driver=vfio-pci no interface", dpdk)
	}

	unbound := byBus["0000:43:00.3"]
	if unbound.Bind != "unbound" || unbound.Driver != "" {
		t.Errorf("unbound device = %+v", unbound)
	}

	// Interface NUMA enrichment still happens for netdev-backed devices.
	if ifaces[0].NUMANode == nil {
		t.Errorf("interface NUMA not enriched")
	}

	// lspci missing is a single non-fatal warning.
	if len(warnings) != 1 || warnings[0].Source != "lspci" || warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one non-fatal lspci warning", warnings)
	}
}

func TestPCIeScanWithLspci(t *testing.T) {
	const bus = "0000:17:00.0"
	fs := pciFS()
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: string(testutil.Fixture(t, "lspci/e810.txt"))},
		"lspci", "-Dnn", "-s", bus,
	)
	ifaces := []model.Interface{{Name: "enp23s0np0", Bus: bus, Driver: "ice"}}

	devices, _ := NewPCIe(fake, fs).Collect(context.Background(), ifaces)
	var d model.PCIDevice
	for _, x := range devices {
		if x.Bus == bus {
			d = x
		}
	}
	if d.Description != "Intel Corporation Ethernet Controller E810-C for QSFP" {
		t.Errorf("description = %q", d.Description)
	}
	if ifaces[0].DeviceName == "" {
		t.Errorf("interface DeviceName not enriched")
	}
}

func TestPCIeScanNoPCITree(t *testing.T) {
	devices, warnings := NewPCIe(run.NewFake(), sysfs.NewFake()).Collect(context.Background(), nil)
	if devices != nil {
		t.Errorf("devices = %v, want nil", devices)
	}
	if len(warnings) != 1 || warnings[0].Source != "sysfs" {
		t.Fatalf("warnings = %v, want one sysfs warning", warnings)
	}
}
