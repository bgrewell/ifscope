package collect

import (
	"sort"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/sysfs"
)

// Bonds collects Linux bonding masters and their members from sysfs. It is
// self-contained (no command runner): bonds are discovered by scanning
// /sys/class/net for a bonding directory, and member/link state comes from each
// device's operstate.
type Bonds struct {
	FS sysfs.FS
}

// NewBonds returns a Bonds collector using fs.
func NewBonds(fs sysfs.FS) *Bonds { return &Bonds{FS: fs} }

// Collect returns the bonds present on the host, sorted by name.
func (c *Bonds) Collect() ([]model.Bond, []model.Warning) {
	entries, err := c.FS.ReadDir("/sys/class/net")
	if err != nil {
		return nil, nil
	}

	var bonds []model.Bond
	for _, e := range entries {
		name := e.Name()
		base := "/sys/class/net/" + name
		if !c.FS.Exists(base + "/bonding/slaves") {
			continue
		}

		b := model.Bond{Name: name, State: c.operstate(name)}
		if mode, err := sysfs.ReadString(c.FS, base+"/bonding/mode"); err == nil {
			b.Mode = friendlyBondMode(firstField(mode))
		}
		if active, err := sysfs.ReadString(c.FS, base+"/bonding/active_slave"); err == nil {
			b.ActiveSlave = active
		}
		if slaves, err := sysfs.ReadString(c.FS, base+"/bonding/slaves"); err == nil {
			for _, m := range strings.Fields(slaves) {
				b.Members = append(b.Members, model.BondMember{Name: m, State: c.operstate(m)})
			}
			sort.Slice(b.Members, func(i, j int) bool { return b.Members[i].Name < b.Members[j].Name })
		}
		bonds = append(bonds, b)
	}
	sort.Slice(bonds, func(i, j int) bool { return bonds[i].Name < bonds[j].Name })
	return bonds, nil
}

// operstate reads an interface's operational state, upper-cased to match the
// interface table (e.g. "UP", "DOWN").
func (c *Bonds) operstate(name string) string {
	s, err := sysfs.ReadString(c.FS, "/sys/class/net/"+name+"/operstate")
	if err != nil {
		return ""
	}
	return strings.ToUpper(s)
}

// firstField returns the first whitespace-separated token of s.
func firstField(s string) string {
	if f := strings.Fields(s); len(f) > 0 {
		return f[0]
	}
	return s
}

// friendlyBondMode annotates the well-known LACP mode; others pass through.
func friendlyBondMode(mode string) string {
	if mode == "802.3ad" {
		return "802.3ad (LACP)"
	}
	return mode
}
