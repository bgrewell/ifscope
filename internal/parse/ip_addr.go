// Package parse holds pure functions that turn raw command output into model
// structs. Parsers perform no I/O, so the fragile bits (JSON shape and text
// scraping) are exhaustively unit testable against fixtures.
package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// ipAddr mirrors one element of `ip -detail -json address show`. Only fields
// ifscope consumes are declared; the rest are ignored.
type ipAddr struct {
	IfIndex   int          `json:"ifindex"`
	IfName    string       `json:"ifname"`
	Flags     []string     `json:"flags"`
	MTU       int          `json:"mtu"`
	OperState string       `json:"operstate"`
	LinkType  string       `json:"link_type"`
	Address   string       `json:"address"`
	Link      string       `json:"link"`
	LinkIndex int          `json:"link_index"`
	Master    string       `json:"master"`
	ParentDev string       `json:"parentdev"`
	ParentBus string       `json:"parentbus"`
	AltNames  []string     `json:"altnames"`
	LinkInfo  *ipLinkInfo  `json:"linkinfo"`
	AddrInfo  []ipAddrInfo `json:"addr_info"`
}

type ipLinkInfo struct {
	InfoKind      string      `json:"info_kind"`
	InfoSlaveKind string      `json:"info_slave_kind"`
	InfoData      *ipInfoData `json:"info_data"`
}

// ipInfoData captures the few info_data fields ifscope uses. Unmarshaling a
// richer info_data (e.g. a bridge's) simply leaves the unused fields zero.
type ipInfoData struct {
	ID       *int   `json:"id"`
	Protocol string `json:"protocol"`
	Type     string `json:"type"`
}

type ipAddrInfo struct {
	Family    string `json:"family"`
	Local     string `json:"local"`
	PrefixLen int    `json:"prefixlen"`
	Scope     string `json:"scope"`
	Dynamic   bool   `json:"dynamic"`
	Secondary bool   `json:"secondary"`
	Label     string `json:"label"`
}

// IPAddresses parses `ip -detail -json address show` into interfaces. VLANs are
// returned alongside physical interfaces; callers partition by Type. Parent
// links are resolved by name and, when only an index is present, by ifindex.
func IPAddresses(data []byte) ([]model.Interface, error) {
	var raw []ipAddr
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse ip address json: %w", err)
	}

	indexToName := make(map[int]string, len(raw))
	for _, e := range raw {
		indexToName[e.IfIndex] = e.IfName
	}

	out := make([]model.Interface, 0, len(raw))
	for _, e := range raw {
		iface := model.Interface{
			ID:         e.IfIndex,
			Name:       e.IfName,
			MAC:        e.Address,
			State:      e.OperState,
			MTU:        e.MTU,
			AltNames:   e.AltNames,
			Bus:        e.ParentDev,
			LinkParent: e.Master,
		}

		for _, a := range e.AddrInfo {
			iface.Addresses = append(iface.Addresses, model.Address{
				Family:    a.Family,
				Local:     a.Local,
				PrefixLen: a.PrefixLen,
				Scope:     a.Scope,
				Dynamic:   a.Dynamic,
				Secondary: a.Secondary,
			})
		}

		iface.Type = classify(e)

		// A VLAN's parent and tag come from link/link_index and info_data.
		if iface.Type == model.TypeVLAN {
			iface.LinkParent = e.Link
			if iface.LinkParent == "" && e.LinkIndex != 0 {
				iface.LinkParent = indexToName[e.LinkIndex]
			}
			if e.LinkInfo != nil && e.LinkInfo.InfoData != nil && e.LinkInfo.InfoData.ID != nil {
				iface.VLANID = *e.LinkInfo.InfoData.ID
			}
		}

		out = append(out, iface)
	}
	return out, nil
}

// classify normalizes an interface into a model.InterfaceType using linkinfo
// kind, link type, and the presence of a PCI parent.
func classify(e ipAddr) model.InterfaceType {
	if e.LinkType == "loopback" {
		return model.TypeLoopback
	}

	if e.LinkInfo != nil && e.LinkInfo.InfoKind != "" {
		switch e.LinkInfo.InfoKind {
		case "vlan":
			return model.TypeVLAN
		case "openvswitch":
			return model.TypeOVS
		case "bridge":
			return model.TypeBridge
		case "bond":
			return model.TypeBond
		case "veth":
			return model.TypeVeth
		case "dummy":
			return model.TypeDummy
		case "tun":
			if e.LinkInfo.InfoData != nil && e.LinkInfo.InfoData.Type == "tap" {
				return model.TypeTap
			}
			return model.TypeTun
		default:
			// A device with a kernel link kind is virtual, not physical. Report
			// the kind verbatim (e.g. macvlan, ipvlan, vxlan, wireguard, vrf,
			// geneve, gre) so the type is accurate rather than forced into a
			// bucket or hidden as unknown.
			return model.InterfaceType(e.LinkInfo.InfoKind)
		}
	}

	// A device backed by a hardware function (PCI, USB, virtio, ...) with no
	// virtual linkinfo is physical. SR-IOV VF refinement (TypeVF) happens later
	// via sysfs (physfn).
	if e.ParentBus != "" || e.ParentDev != "" {
		return model.TypePhysical
	}

	// No link kind and no device backing: classification is undetermined rather
	// than assumed physical, to avoid reporting inaccurate information.
	return model.TypeUnknown
}
