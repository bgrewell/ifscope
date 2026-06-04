package collect

import (
	"sort"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// Bridges collects native Linux bridges and their member ports from sysfs.
// (Open vSwitch bridges are collected separately by the OVS collector.)
type Bridges struct {
	FS sysfs.FS
}

// NewBridges returns a Bridges collector using fs.
func NewBridges(fs sysfs.FS) *Bridges { return &Bridges{FS: fs} }

// Collect returns the Linux bridges present on the host, sorted by name. A
// bridge is any /sys/class/net entry with a bridge/ directory; its members are
// the entries under brif/.
func (c *Bridges) Collect() ([]model.Bridge, []model.Warning) {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil, nil
	}

	var bridges []model.Bridge
	for _, e := range entries {
		name := e.Name()
		base := "/sys/class/net/" + name
		if !c.FS.Exists(base + "/bridge") {
			continue
		}

		b := model.Bridge{Name: name, State: linkOperstate(c.FS, name)}
		if stp, err := sysfs.ReadInt(c.FS, base+"/bridge/stp_state"); err == nil && stp > 0 {
			b.STP = true
		}
		if members, err := c.FS.ReadDir(base + "/brif"); err == nil {
			for _, m := range members {
				b.Members = append(b.Members, model.BridgeMember{
					Name:  m.Name(),
					State: linkOperstate(c.FS, m.Name()),
				})
			}
			sort.Slice(b.Members, func(i, j int) bool { return b.Members[i].Name < b.Members[j].Name })
		}
		bridges = append(bridges, b)
	}
	sort.Slice(bridges, func(i, j int) bool { return bridges[i].Name < bridges[j].Name })
	return bridges, nil
}

// linkOperstate reads an interface's operational state from sysfs, upper-cased
// to match the interface table (e.g. "UP", "DOWN").
func linkOperstate(fs sysfs.FS, name string) string {
	s, err := sysfs.ReadString(fs, "/sys/class/net/"+name+"/operstate")
	if err != nil {
		return ""
	}
	return strings.ToUpper(s)
}
