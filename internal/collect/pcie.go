package collect

import (
	"context"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

const lspciCmd = "lspci"

// PCIe builds the PCIe device table for interfaces backed by a PCI function,
// correlating lspci descriptions with sysfs attributes.
type PCIe struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewPCIe returns a PCIe collector using r and fs.
func NewPCIe(r run.Runner, fs sysfs.FS) *PCIe { return &PCIe{Runner: r, FS: fs} }

// Collect returns one PCIDevice per PCI-backed interface and, as a side effect,
// fills each interface's NUMA node and device description. A missing lspci is a
// non-fatal warning; sysfs gaps simply leave fields empty.
func (c *PCIe) Collect(ctx context.Context, ifaces []model.Interface) ([]model.PCIDevice, []model.Warning) {
	var devices []model.PCIDevice
	var warnings []model.Warning
	lspciMissing := false

	for i := range ifaces {
		bus := ifaces[i].Bus
		if !isPCIBus(bus) {
			continue
		}

		dev := model.PCIDevice{Bus: bus, Interface: ifaces[i].Name, Driver: ifaces[i].Driver}
		base := "/sys/bus/pci/devices/" + bus
		dev.VendorID = readHexID(c.FS, base+"/vendor")
		dev.DeviceID = readHexID(c.FS, base+"/device")
		dev.SubsystemVendorID = readHexID(c.FS, base+"/subsystem_vendor")
		dev.SubsystemDeviceID = readHexID(c.FS, base+"/subsystem_device")
		if n, err := sysfs.ReadInt(c.FS, base+"/numa_node"); err == nil && n >= 0 {
			node := n
			dev.NUMANode = &node
			ifaces[i].NUMANode = &node
		}
		if v, err := sysfs.ReadString(c.FS, base+"/current_link_speed"); err == nil {
			dev.LinkSpeed = v
		}
		if v, err := sysfs.ReadString(c.FS, base+"/current_link_width"); err == nil {
			dev.LinkWidth = v
		}

		if !lspciMissing {
			out, _, err := c.Runner.Run(ctx, lspciCmd, "-Dnn", "-s", bus)
			switch {
			case err == nil:
				info := parse.LspciDevice(out)
				dev.Description = info.Description
				if dev.VendorID == "" {
					dev.VendorID = info.VendorID
				}
				if dev.DeviceID == "" {
					dev.DeviceID = info.DeviceID
				}
				ifaces[i].DeviceName = info.Description
			case run.IsNotFound(err):
				lspciMissing = true
			}
		}

		devices = append(devices, dev)
	}

	if lspciMissing {
		warnings = append(warnings, model.Warning{
			Source:  "lspci",
			Message: "lspci not found; PCIe device descriptions unavailable",
		})
	}
	return devices, warnings
}

// readHexID reads a sysfs hex id ("0x8086") and returns it without the prefix.
func readHexID(fs sysfs.FS, path string) string {
	v, err := sysfs.ReadString(fs, path)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(v, "0x")
}
