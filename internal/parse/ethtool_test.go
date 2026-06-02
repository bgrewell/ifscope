package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestEthtoolDriverInfo(t *testing.T) {
	got := EthtoolDriverInfo(testutil.Fixture(t, "ethtool/ice-i.txt"))
	if got.Driver != "ice" {
		t.Errorf("driver = %q, want ice", got.Driver)
	}
	if got.Firmware != "4.30 0x8001be17 1.3429.0" {
		t.Errorf("firmware = %q", got.Firmware)
	}
	if got.Bus != "0000:17:00.0" {
		t.Errorf("bus = %q", got.Bus)
	}
}

func TestEthtoolDriverInfoVirtualNAFirmware(t *testing.T) {
	got := EthtoolDriverInfo(testutil.Fixture(t, "ethtool/virtual-i.txt"))
	if got.Driver != "bridge" {
		t.Errorf("driver = %q, want bridge", got.Driver)
	}
	if got.Firmware != "" {
		t.Errorf("firmware = %q, want empty (N/A normalized)", got.Firmware)
	}
}

func TestEthtoolSettings(t *testing.T) {
	got := EthtoolSettings(testutil.Fixture(t, "ethtool/ice-settings.txt"))
	if got.Speed != "100 Gb/s" {
		t.Errorf("speed = %q, want 100 Gb/s", got.Speed)
	}
	if got.Port != "Direct Attach Copper" {
		t.Errorf("port = %q", got.Port)
	}
	if got.Duplex != "Full" {
		t.Errorf("duplex = %q", got.Duplex)
	}
}

func TestNormalizeSpeed(t *testing.T) {
	cases := map[string]string{
		"100000Mb/s": "100 Gb/s",
		"25000Mb/s":  "25 Gb/s",
		"1000Mb/s":   "1 Gb/s",
		"100Mb/s":    "100 Mb/s",
		"Unknown!":   "",
		"":           "",
	}
	for in, want := range cases {
		if got := normalizeSpeed(in); got != want {
			t.Errorf("normalizeSpeed(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLspciDevice(t *testing.T) {
	got := LspciDevice(testutil.Fixture(t, "lspci/e810.txt"))
	if got.VendorID != "8086" || got.DeviceID != "1592" {
		t.Errorf("ids = %s:%s, want 8086:1592", got.VendorID, got.DeviceID)
	}
	want := "Intel Corporation Ethernet Controller E810-C for QSFP"
	if got.Description != want {
		t.Errorf("description = %q, want %q", got.Description, want)
	}
}

func TestLspciDeviceEmpty(t *testing.T) {
	if got := LspciDevice(nil); got != (PCIInfo{}) {
		t.Errorf("empty input = %+v, want zero", got)
	}
}
