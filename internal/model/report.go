package model

// Host carries basic identity of the inspected host.
type Host struct {
	Hostname string `json:"hostname,omitempty"`
}

// Report is the top-level correlated result produced by a run. Renderers emit
// it as a human table view or as JSON. The Version field versions the JSON
// schema so consumers can adapt to future changes.
type Report struct {
	Version     string            `json:"version"`
	Host        Host              `json:"host"`
	Interfaces  []Interface       `json:"interfaces,omitempty"`
	VLANs       []Interface       `json:"vlans,omitempty"`
	Bonds       []Bond            `json:"bonds,omitempty"`
	Bridges     []Bridge          `json:"bridges,omitempty"`
	BridgeVLANs []BridgeVLAN      `json:"bridge_vlans,omitempty"`
	Tunnels     []Tunnel          `json:"tunnels,omitempty"`
	WireGuard   []WireGuardDevice `json:"wireguard,omitempty"`
	Routes      []Route           `json:"routes,omitempty"`
	Rules       []Rule            `json:"rules,omitempty"`
	Neighbors   []Neighbor        `json:"neighbors,omitempty"`
	LLDP        []LLDPNeighbor    `json:"lldp,omitempty"`
	FDB         []FDBEntry        `json:"fdb,omitempty"`
	DNS         []DNS             `json:"dns,omitempty"`
	PCIe        []PCIDevice       `json:"pcie,omitempty"`
	Devlink     []DevlinkPort     `json:"devlink,omitempty"`
	Qdiscs      []Qdisc           `json:"qdiscs,omitempty"`
	Offloads    []Offloads        `json:"offloads,omitempty"`
	Queues      []Queues          `json:"queues,omitempty"`
	IRQs        []IRQ             `json:"irqs,omitempty"`
	Multicast   []MulticastGroup  `json:"multicast,omitempty"`
	OVS         *OVS              `json:"ovs,omitempty"`
	Netns       []Netns           `json:"netns,omitempty"`
	Stats       []InterfaceStats  `json:"stats,omitempty"`
	Sockets     []Socket          `json:"sockets,omitempty"`
	Tests       []TestResult      `json:"tests,omitempty"`
	Warnings    []Warning         `json:"warnings,omitempty"`
}

// AddWarning appends a warning to the report.
func (r *Report) AddWarning(w Warning) {
	r.Warnings = append(r.Warnings, w)
}
