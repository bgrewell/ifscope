// Package correlate stitches collector outputs into a final model.Report:
// partitioning interfaces from VLANs, applying filters, and sorting for stable
// presentation.
package correlate

import (
	"sort"

	"github.com/bgrewell/ifscope/internal/model"
)

// Partition splits a flat interface list into non-VLAN interfaces and VLANs.
func Partition(ifaces []model.Interface) (interfaces, vlans []model.Interface) {
	for _, i := range ifaces {
		if i.Type == model.TypeVLAN {
			vlans = append(vlans, i)
		} else {
			interfaces = append(interfaces, i)
		}
	}
	return interfaces, vlans
}

// stateRank orders interfaces UP first, then DOWN, then everything else.
func stateRank(state string) int {
	switch state {
	case "UP":
		return 0
	case "DOWN":
		return 1
	default:
		return 2
	}
}

// SortInterfaces sorts in place by state rank then name.
func SortInterfaces(ifaces []model.Interface) {
	sort.SliceStable(ifaces, func(a, b int) bool {
		ra, rb := stateRank(ifaces[a].State), stateRank(ifaces[b].State)
		if ra != rb {
			return ra < rb
		}
		return ifaces[a].Name < ifaces[b].Name
	})
}

// SortVLANs sorts in place by parent device then VLAN id.
func SortVLANs(vlans []model.Interface) {
	sort.SliceStable(vlans, func(a, b int) bool {
		if vlans[a].LinkParent != vlans[b].LinkParent {
			return vlans[a].LinkParent < vlans[b].LinkParent
		}
		return vlans[a].VLANID < vlans[b].VLANID
	})
}

// IsPhysical reports whether an interface is a hardware function (a physical
// port or an SR-IOV virtual function).
func IsPhysical(i model.Interface) bool {
	return i.Type == model.TypePhysical || i.Type == model.TypeVF
}

// Filter selects interfaces by common criteria. The zero value matches all.
type Filter struct {
	Up       bool
	Name     string
	Driver   string
	State    string
	Physical bool
	Virtual  bool
	VF       bool
	PF       bool
}

// Match reports whether i satisfies every active filter criterion.
func (f Filter) Match(i model.Interface) bool {
	if f.Up && i.State != "UP" {
		return false
	}
	if f.Name != "" && i.Name != f.Name {
		return false
	}
	if f.Driver != "" && i.Driver != f.Driver {
		return false
	}
	if f.State != "" && !equalFold(i.State, f.State) {
		return false
	}
	if f.Physical && !IsPhysical(i) {
		return false
	}
	if f.Virtual && IsPhysical(i) {
		return false
	}
	if f.VF && (i.SRIOV == nil || !i.SRIOV.VF) {
		return false
	}
	if f.PF && (i.SRIOV == nil || !i.SRIOV.Capable || i.SRIOV.VF) {
		return false
	}
	return true
}

// Apply returns the subset of ifaces matching f, preserving order.
func (f Filter) Apply(ifaces []model.Interface) []model.Interface {
	out := make([]model.Interface, 0, len(ifaces))
	for _, i := range ifaces {
		if f.Match(i) {
			out = append(out, i)
		}
	}
	return out
}

// equalFold compares ASCII state strings case-insensitively without allocating
// for the common already-matching case.
func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if 'A' <= ca && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if 'A' <= cb && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
