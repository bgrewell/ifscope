package model

// SRIOVInfo describes the SR-IOV state of an interface, whether it is a
// physical function (PF) exposing virtual functions, or itself a VF.
//
// Invariant: Enabled is true exactly when ConfiguredVFs > 0.
type SRIOVInfo struct {
	Capable       bool   `json:"capable"`
	TotalVFs      int    `json:"total_vfs,omitempty"`
	ConfiguredVFs int    `json:"configured_vfs,omitempty"`
	Enabled       bool   `json:"enabled,omitempty"`
	VF            bool   `json:"vf"`
	PF            string `json:"pf,omitempty"`
	PFBus         string `json:"pf_bus,omitempty"`
	VFIndex       *int   `json:"vf_index,omitempty"`
	VFs           []VF   `json:"vfs,omitempty"`
}

// VF describes a single virtual function exposed by a PF.
type VF struct {
	Index      int    `json:"index"`
	Bus        string `json:"bus,omitempty"`
	Netdev     string `json:"netdev,omitempty"`
	Driver     string `json:"driver,omitempty"`
	MAC        string `json:"mac,omitempty"`
	VLAN       int    `json:"vlan,omitempty"`
	SpoofCheck bool   `json:"spoof_check,omitempty"`
	Trust      bool   `json:"trust,omitempty"`
	LinkState  string `json:"link_state,omitempty"`
}
