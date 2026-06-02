package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestPCIeCollect(t *testing.T) {
	const bus = "0000:17:00.0"
	fs := sysfs.NewFake()
	base := "/sys/bus/pci/devices/" + bus
	fs.Files[base+"/vendor"] = "0x8086\n"
	fs.Files[base+"/device"] = "0x1592\n"
	fs.Files[base+"/numa_node"] = "1\n"
	fs.Files[base+"/current_link_speed"] = "16.0 GT/s PCIe\n"
	fs.Files[base+"/current_link_width"] = "8\n"

	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "lspci/e810.txt"))}, "lspci", "-Dnn", "-s", bus)

	ifaces := []model.Interface{
		{Name: "eth0", Type: model.TypePhysical, Bus: bus, Driver: "ice"},
		{Name: "br0", Type: model.TypeBridge}, // no PCI bus; skipped
	}
	devices, warnings := NewPCIe(fake, fs).Collect(context.Background(), ifaces)
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
	if len(devices) != 1 {
		t.Fatalf("devices = %d, want 1", len(devices))
	}

	d := devices[0]
	if d.Interface != "eth0" || d.Driver != "ice" {
		t.Errorf("device interface/driver = %q/%q", d.Interface, d.Driver)
	}
	if d.VendorID != "8086" || d.DeviceID != "1592" {
		t.Errorf("ids = %s:%s, want 8086:1592", d.VendorID, d.DeviceID)
	}
	if d.Description != "Intel Corporation Ethernet Controller E810-C for QSFP" {
		t.Errorf("description = %q", d.Description)
	}
	if d.NUMANode == nil || *d.NUMANode != 1 {
		t.Errorf("numa = %v, want 1", d.NUMANode)
	}
	if ifaces[0].NUMANode == nil || *ifaces[0].NUMANode != 1 {
		t.Errorf("interface NUMA not enriched: %v", ifaces[0].NUMANode)
	}
	if ifaces[0].DeviceName == "" {
		t.Errorf("interface DeviceName not enriched")
	}
}

func TestPCIeCollectMissingLspci(t *testing.T) {
	const bus = "0000:17:00.0"
	fs := sysfs.NewFake()
	fs.Files["/sys/bus/pci/devices/"+bus+"/numa_node"] = "-1\n"

	ifaces := []model.Interface{{Name: "eth0", Bus: bus}}
	devices, warnings := NewPCIe(run.NewFake(), fs).Collect(context.Background(), ifaces)

	if len(devices) != 1 {
		t.Fatalf("devices = %d, want 1", len(devices))
	}
	if devices[0].NUMANode != nil {
		t.Errorf("numa_node -1 should be treated as unknown, got %v", devices[0].NUMANode)
	}
	if len(warnings) != 1 || warnings[0].Source != "lspci" || warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one non-fatal lspci warning", warnings)
	}
}
