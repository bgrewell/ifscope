package collect

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// IRQ collects NIC interrupt CPU affinity from /proc and /sys (no subprocess).
type IRQ struct {
	FS sysfs.FS
}

// NewIRQ returns an IRQ collector using fs.
func NewIRQ(fs sysfs.FS) *IRQ { return &IRQ{FS: fs} }

// Collect parses /proc/interrupts, keeps only IRQs whose name maps to a known
// interface, and reads each one's SMP affinity.
func (c *IRQ) Collect() ([]model.IRQ, []model.Warning) {
	data, err := c.FS.ReadFile("/proc/interrupts")
	if err != nil {
		return nil, []model.Warning{{Source: "irq", Message: "cannot read /proc/interrupts; IRQ data unavailable"}}
	}

	names := c.interfaceNames()
	var out []model.IRQ
	for _, line := range parse.Interrupts(data) {
		dev := matchDevice(line.Name, names)
		if dev == "" {
			continue
		}
		cpus, _ := sysfs.ReadString(c.FS, fmt.Sprintf("/proc/irq/%d/smp_affinity_list", line.Number))
		out = append(out, model.IRQ{Device: dev, IRQ: line.Number, Name: line.Name, CPUs: cpus})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Device != out[j].Device {
			return out[i].Device < out[j].Device
		}
		return out[i].IRQ < out[j].IRQ
	})
	return out, nil
}

// interfaceNames lists netdev names, longest first so device matching prefers
// the most specific name.
func (c *IRQ) interfaceNames() []string {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.Name() != "lo" {
			names = append(names, e.Name())
		}
	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) > len(names[j]) })
	return names
}

// matchDevice returns the first (longest) interface name that appears in the
// IRQ name, or "".
func matchDevice(irqName string, names []string) string {
	for _, n := range names {
		if strings.Contains(irqName, n) {
			return n
		}
	}
	return ""
}
