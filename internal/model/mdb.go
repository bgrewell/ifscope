package model

// MDBEntry is a bridge multicast database entry: a multicast group forwarded to
// a specific bridge port.
type MDBEntry struct {
	Bridge string `json:"bridge"`
	Port   string `json:"port"`
	Group  string `json:"group"`
	VLAN   int    `json:"vlan,omitempty"`
	State  string `json:"state,omitempty"`
}
