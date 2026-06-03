package model

// InterfaceType is a normalized classification of a network interface. For
// virtual devices without a dedicated constant below (e.g. macvlan, ipvlan,
// vxlan, wireguard, vrf, geneve, gre), the value is the kernel's link kind
// verbatim, so the type is always accurate. Unknown is used only when the type
// genuinely cannot be determined.
type InterfaceType string

// Recognized interface types.
const (
	TypePhysical InterfaceType = "physical"
	TypeVLAN     InterfaceType = "vlan"
	TypeBond     InterfaceType = "bond"
	TypeBridge   InterfaceType = "bridge"
	TypeOVS      InterfaceType = "ovs"
	TypeTap      InterfaceType = "tap"
	TypeTun      InterfaceType = "tun"
	TypeVeth     InterfaceType = "veth"
	TypeDummy    InterfaceType = "dummy"
	TypeLoopback InterfaceType = "loopback"
	TypeVF       InterfaceType = "vf"
	TypeUnknown  InterfaceType = "unknown"
)

// Interface is the correlated view of a single network interface, stitching
// together data from `ip`, ethtool, sysfs, lspci, OVS, and SR-IOV sources.
//
// VLAN interfaces are represented as Interface values with Type == TypeVLAN;
// for them LinkParent holds the parent device and VLANID holds the tag.
type Interface struct {
	ID         int           `json:"id"`
	Name       string        `json:"name"`
	Type       InterfaceType `json:"type"`
	MAC        string        `json:"mac,omitempty"`
	State      string        `json:"state"`
	MTU        int           `json:"mtu,omitempty"`
	Addresses  []Address     `json:"addresses,omitempty"`
	LinkParent string        `json:"link_parent,omitempty"`
	AltNames   []string      `json:"altnames,omitempty"`
	VLANID     int           `json:"vlan_id,omitempty"`

	Driver     string `json:"driver,omitempty"`
	Firmware   string `json:"firmware,omitempty"`
	Bus        string `json:"bus,omitempty"`
	Speed      string `json:"speed,omitempty"`
	Port       string `json:"port,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	NUMANode   *int   `json:"numa_node,omitempty"`

	SRIOV *SRIOVInfo     `json:"sriov,omitempty"`
	OVS   *OVSMembership `json:"ovs,omitempty"`
}

// IPv4 returns the interface's IPv4 addresses.
func (i Interface) IPv4() []Address {
	var out []Address
	for _, a := range i.Addresses {
		if a.IsIPv4() {
			out = append(out, a)
		}
	}
	return out
}
