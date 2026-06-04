package model

// Neighbor is an entry from the ARP (IPv4) / NDP (IPv6) neighbor table.
type Neighbor struct {
	Dst    string `json:"dst"`
	Dev    string `json:"dev,omitempty"`
	LLAddr string `json:"lladdr,omitempty"`
	State  string `json:"state,omitempty"`
	Router bool   `json:"router,omitempty"`
	Family string `json:"family,omitempty"`
}
