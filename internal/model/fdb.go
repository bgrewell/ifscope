package model

// FDBEntry is one entry of a bridge forwarding database (learned or static MAC
// to port mapping).
type FDBEntry struct {
	MAC    string   `json:"mac"`
	Dev    string   `json:"dev"`
	VLAN   *int     `json:"vlan,omitempty"`
	Master string   `json:"master,omitempty"`
	State  string   `json:"state,omitempty"`
	Flags  []string `json:"flags,omitempty"`
}
