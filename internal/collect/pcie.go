package collect

import (
	"context"
	"path"
	"sort"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

const (
	lspciCmd    = "lspci"
	pciDevices  = "/sys/bus/pci/devices"
	ethClassHex = "0x0200" // PCI class prefix for Ethernet controllers
)

// dpdkDrivers are userspace/passthrough drivers used to detach a NIC from the
// kernel (e.g. for DPDK). A device bound to one of these has no kernel netdev.
var dpdkDrivers = map[string]bool{
	"vfio-pci":        true,
	"uio_pci_generic": true,
	"igb_uio":         true,
	"uio":             true,
}

// PCIe builds the PCIe network-device table by scanning sysfs for Ethernet-class
// devices, independent of kernel netdevs. This surfaces NICs detached from the
// kernel (DPDK/passthrough or unbound) alongside ordinary netdev-backed ones,
// correlating lspci descriptions and interface names where available.
type PCIe struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewPCIe returns a PCIe collector using r and fs.
func NewPCIe(r run.Runner, fs sysfs.FS) *PCIe { return &PCIe{Runner: r, FS: fs} }

// Collect enumerates Ethernet-class PCI devices. For netdev-backed devices it
// also fills the owning interface's NUMA node and device description. A missing
// lspci or an unreadable PCI tree degrades to a non-fatal warning.
func (c *PCIe) Collect(ctx context.Context, ifaces []model.Interface) ([]model.PCIDevice, []model.Warning) {
	busToIdx := map[string]int{}
	for i := range ifaces {
		if ifaces[i].Bus != "" {
			busToIdx[ifaces[i].Bus] = i
		}
	}

	entries, err := c.FS.ReadDir(pciDevices)
	if err != nil {
		return nil, []model.Warning{{Source: "sysfs", Message: "cannot scan " + pciDevices + "; PCIe data unavailable"}}
	}

	var devices []model.PCIDevice
	var warnings []model.Warning
	lspciMissing := false

	for _, e := range entries {
		bus := e.Name()
		base := pciDevices + "/" + bus
		if class, _ := sysfs.ReadString(c.FS, base+"/class"); !strings.HasPrefix(class, ethClassHex) {
			continue
		}

		dev := model.PCIDevice{Bus: bus}
		if t, err := c.FS.Readlink(base + "/driver"); err == nil {
			dev.Driver = path.Base(t)
		}

		// Resolve the kernel netdev: prefer the interface list, else the
		// device's own net/ directory.
		idx, hasIface := busToIdx[bus]
		if hasIface {
			dev.Interface = ifaces[idx].Name
		} else {
			dev.Interface = c.firstNetdev(base)
		}
		dev.Bind = classifyBind(dev.Interface, dev.Driver)

		dev.VendorID = readHexID(c.FS, base+"/vendor")
		dev.DeviceID = readHexID(c.FS, base+"/device")
		dev.SubsystemVendorID = readHexID(c.FS, base+"/subsystem_vendor")
		dev.SubsystemDeviceID = readHexID(c.FS, base+"/subsystem_device")
		if n, err := sysfs.ReadInt(c.FS, base+"/numa_node"); err == nil && n >= 0 {
			node := n
			dev.NUMANode = &node
			if hasIface {
				ifaces[idx].NUMANode = &node
			}
		}
		if v, err := sysfs.ReadString(c.FS, base+"/current_link_speed"); err == nil {
			dev.LinkSpeed = v
		}
		if v, err := sysfs.ReadString(c.FS, base+"/current_link_width"); err == nil {
			dev.LinkWidth = v
		}

		if !lspciMissing {
			out, _, lerr := c.Runner.Run(ctx, lspciCmd, "-Dnn", "-s", bus)
			switch {
			case lerr == nil:
				info := parse.LspciDevice(out)
				dev.Description = info.Description
				if dev.VendorID == "" {
					dev.VendorID = info.VendorID
				}
				if dev.DeviceID == "" {
					dev.DeviceID = info.DeviceID
				}
				if hasIface {
					ifaces[idx].DeviceName = info.Description
				}
			case run.IsNotFound(lerr):
				lspciMissing = true
			}
		}

		devices = append(devices, dev)
	}

	sort.Slice(devices, func(i, j int) bool { return devices[i].Bus < devices[j].Bus })

	if lspciMissing {
		warnings = append(warnings, model.Warning{
			Source:  "lspci",
			Message: "lspci not found; PCIe device descriptions unavailable",
		})
	}
	return devices, warnings
}

// firstNetdev returns the first interface name under a PCI device's net/
// directory, or "" when the device exposes no kernel netdev.
func (c *PCIe) firstNetdev(base string) string {
	entries, err := c.FS.ReadDir(base + "/net")
	if err != nil || len(entries) == 0 {
		return ""
	}
	return entries[0].Name()
}

// classifyBind reports how a PCI device is attached to the kernel.
func classifyBind(netdev, driver string) string {
	switch {
	case netdev != "":
		return "kernel"
	case dpdkDrivers[driver]:
		return "dpdk"
	case driver == "":
		return "unbound"
	default:
		return "detached"
	}
}

// readHexID reads a sysfs hex id ("0x8086") and returns it without the prefix.
func readHexID(fs sysfs.FS, path string) string {
	v, err := sysfs.ReadString(fs, path)
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(v, "0x")
}
