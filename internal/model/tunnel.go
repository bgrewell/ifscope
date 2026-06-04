package model

// Tunnel is an overlay/tunnel interface (VXLAN, GENEVE, GRE, IPIP, SIT, …).
// VNI and Port are pointers so "absent" is distinguishable from zero.
type Tunnel struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	VNI    *int   `json:"vni,omitempty"`
	Local  string `json:"local,omitempty"`
	Remote string `json:"remote,omitempty"`
	Port   *int   `json:"port,omitempty"`
	TTL    string `json:"ttl,omitempty"`
	Dev    string `json:"dev,omitempty"`
}
