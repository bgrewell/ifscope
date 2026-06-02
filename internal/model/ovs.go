package model

// OVS captures the Open vSwitch topology: bridges, ports, and interfaces.
type OVS struct {
	Bridges    []OVSBridge    `json:"bridges,omitempty"`
	Ports      []OVSPort      `json:"ports,omitempty"`
	Interfaces []OVSInterface `json:"interfaces,omitempty"`
}

// OVSBridge is an Open vSwitch bridge and the ports attached to it.
type OVSBridge struct {
	Name  string   `json:"name"`
	Ports []string `json:"ports,omitempty"`
}

// OVSPort is a port on a bridge. Tag holds the access VLAN when set; Trunks
// holds the allowed VLANs in trunk mode.
type OVSPort struct {
	Name       string   `json:"name"`
	Bridge     string   `json:"bridge,omitempty"`
	Interfaces []string `json:"interfaces,omitempty"`
	Tag        *int     `json:"tag,omitempty"`
	Trunks     []int    `json:"trunks,omitempty"`
	VLANMode   string   `json:"vlan_mode,omitempty"`
	BondMode   string   `json:"bond_mode,omitempty"`
	LACP       string   `json:"lacp,omitempty"`
}

// OVSInterface is an interface record within Open vSwitch.
type OVSInterface struct {
	Name        string            `json:"name"`
	Type        string            `json:"type,omitempty"`
	OFPort      *int              `json:"ofport,omitempty"`
	ExternalIDs map[string]string `json:"external_ids,omitempty"`
	Options     map[string]string `json:"options,omitempty"`
}

// OVSMembership annotates a Linux interface with its OVS bridge/port context.
type OVSMembership struct {
	Bridge string `json:"bridge,omitempty"`
	Port   string `json:"port,omitempty"`
	Tag    *int   `json:"tag,omitempty"`
}
