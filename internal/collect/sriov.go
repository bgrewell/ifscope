package collect

import (
	"context"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// SRIOV enriches interfaces with SR-IOV physical-function and virtual-function
// state, read primarily from sysfs and supplemented with VF attributes from
// `ip -details -json link`.
type SRIOV struct {
	Runner run.Runner
	FS     sysfs.FS
}

// NewSRIOV returns an SR-IOV collector using r and fs.
func NewSRIOV(r run.Runner, fs sysfs.FS) *SRIOV { return &SRIOV{Runner: r, FS: fs} }

// Enrich fills each interface's SRIOV field in place and reclassifies VFs as
// TypeVF. A bus→netdev map (built from the interface list) resolves VF and PF
// netdev names without extra sysfs walks.
func (c *SRIOV) Enrich(ctx context.Context, ifaces []model.Interface) []model.Warning {
	busToName := map[string]string{}
	for _, i := range ifaces {
		if i.Bus != "" {
			busToName[i.Bus] = i.Name
		}
	}

	for i := range ifaces {
		dev := "/sys/class/net/" + ifaces[i].Name + "/device"
		switch {
		case c.FS.Exists(dev + "/sriov_totalvfs"):
			ifaces[i].SRIOV = c.collectPF(ctx, ifaces[i], dev, busToName)
		case c.FS.Exists(dev + "/physfn"):
			ifaces[i].SRIOV = c.collectVF(ifaces[i], dev, busToName)
			ifaces[i].Type = model.TypeVF
		}
	}
	return nil
}

// collectPF reads PF SR-IOV state and enumerates its virtual functions.
func (c *SRIOV) collectPF(ctx context.Context, iface model.Interface, dev string, busToName map[string]string) *model.SRIOVInfo {
	total, _ := sysfs.ReadInt(c.FS, dev+"/sriov_totalvfs")
	configured, _ := sysfs.ReadInt(c.FS, dev+"/sriov_numvfs")

	// sriov_totalvfs is the read-only hardware/firmware maximum. A value of 0
	// means no VFs are possible, so the device is not SR-IOV capable in
	// practice; report it as an ordinary interface (no SR-IOV block).
	if total <= 0 {
		return nil
	}

	info := &model.SRIOVInfo{
		Capable:       true,
		TotalVFs:      total,
		ConfiguredVFs: configured,
		Enabled:       configured > 0,
	}
	if configured == 0 {
		return info
	}

	// kernel-reported VF attributes (mac/vlan/spoof/trust/link-state), if any.
	attrs := c.vfAttrs(ctx, iface.Name)

	entries, _ := c.FS.ReadDir(dev)
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "virtfn") {
			continue
		}
		idx := atoiDefault(strings.TrimPrefix(name, "virtfn"), -1)
		if idx < 0 {
			continue
		}
		target, err := c.FS.Readlink(dev + "/" + name)
		if err != nil {
			continue
		}
		vfBus := path.Base(target)

		vf := model.VF{
			Index:  idx,
			Bus:    vfBus,
			Netdev: busToName[vfBus],
			Driver: c.driverOf(vfBus),
		}
		if a, ok := attrs[idx]; ok {
			vf.MAC = a.MAC
			vf.VLAN = a.VLAN
			vf.SpoofCheck = a.SpoofCheck
			vf.Trust = a.Trust
			vf.LinkState = a.LinkState
		}
		info.VFs = append(info.VFs, vf)
	}
	sortVFs(info.VFs)
	return info
}

// collectVF reads VF state: its PF, PF bus, and VF index.
func (c *SRIOV) collectVF(iface model.Interface, dev string, busToName map[string]string) *model.SRIOVInfo {
	info := &model.SRIOVInfo{VF: true}

	if target, err := c.FS.Readlink(dev + "/physfn"); err == nil {
		pfBus := path.Base(target)
		info.PFBus = pfBus
		info.PF = busToName[pfBus]
		if idx := c.vfIndex(pfBus, iface.Bus); idx >= 0 {
			info.VFIndex = &idx
		}
	}
	return info
}

// vfIndex finds which virtfnN of the PF maps to vfBus, returning the index.
func (c *SRIOV) vfIndex(pfBus, vfBus string) int {
	dev := "/sys/bus/pci/devices/" + pfBus
	entries, err := c.FS.ReadDir(dev)
	if err != nil {
		return -1
	}
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, "virtfn") {
			continue
		}
		if target, err := c.FS.Readlink(dev + "/" + name); err == nil && path.Base(target) == vfBus {
			return atoiDefault(strings.TrimPrefix(name, "virtfn"), -1)
		}
	}
	return -1
}

// driverOf returns the driver bound to a PCI function (kernel module name,
// "vfio-pci", or "" when unbound).
func (c *SRIOV) driverOf(bus string) string {
	target, err := c.FS.Readlink("/sys/bus/pci/devices/" + bus + "/driver")
	if err != nil {
		return ""
	}
	return path.Base(target)
}

// vfAttrs fetches VF attributes for a PF via `ip -details -json link show`.
func (c *SRIOV) vfAttrs(ctx context.Context, pf string) map[int]model.VF {
	out, _, err := c.Runner.Run(ctx, ipCmd, "-details", "-json", "link", "show", pf)
	if err != nil {
		return nil
	}
	attrs, perr := parse.VFAttrs(out)
	if perr != nil {
		return nil
	}
	return attrs
}

// atoiDefault parses s as an int, returning def on failure.
func atoiDefault(s string, def int) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// sortVFs orders virtual functions by index.
func sortVFs(vfs []model.VF) {
	sort.Slice(vfs, func(a, b int) bool { return vfs[a].Index < vfs[b].Index })
}
