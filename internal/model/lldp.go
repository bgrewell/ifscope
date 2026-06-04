package model

// LLDPNeighbor is a link-layer (LLDP) neighbor discovered on a local port.
type LLDPNeighbor struct {
	LocalPort    string   `json:"local_port"`
	Chassis      string   `json:"chassis,omitempty"`
	ChassisID    string   `json:"chassis_id,omitempty"`
	ChassisDescr string   `json:"chassis_descr,omitempty"`
	MgmtIPs      []string `json:"mgmt_ips,omitempty"`
	PortID       string   `json:"port_id,omitempty"`
	PortDescr    string   `json:"port_descr,omitempty"`
	TTL          string   `json:"ttl,omitempty"`
	Capabilities []string `json:"capabilities,omitempty"`
}
