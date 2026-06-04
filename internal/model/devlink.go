package model

// DevlinkPort is a devlink port — the device-level view of a NIC port,
// including its SR-IOV flavour (physical function, PF, or VF). This complements
// the SR-IOV view for switchdev/offload setups.
type DevlinkPort struct {
	Handle  string `json:"handle"`
	Type    string `json:"type,omitempty"`
	Flavour string `json:"flavour,omitempty"`
	Netdev  string `json:"netdev,omitempty"`
	PfNum   *int   `json:"pfnum,omitempty"`
	VfNum   *int   `json:"vfnum,omitempty"`
	Lanes   int    `json:"lanes,omitempty"`
}
