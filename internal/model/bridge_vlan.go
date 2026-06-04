package model

// BridgeVLAN is a VLAN (or VLAN range) configured on a bridge port via the
// kernel's VLAN filtering. PVID/untagged status is conveyed in Flags.
type BridgeVLAN struct {
	Port    string   `json:"port"`
	VLAN    int      `json:"vlan"`
	VLANEnd int      `json:"vlan_end,omitempty"`
	Flags   []string `json:"flags,omitempty"`
}
