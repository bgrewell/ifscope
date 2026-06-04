package model

// TCClass is a traffic-control class (e.g. an htb/hfsc shaping class). Rate and
// Ceil are in bits per second (tc reports bytes/s; ifscope normalizes to bits).
type TCClass struct {
	Dev    string `json:"dev"`
	Kind   string `json:"kind"`
	Handle string `json:"handle"`
	Parent string `json:"parent,omitempty"`
	Root   bool   `json:"root,omitempty"`
	Rate   uint64 `json:"rate,omitempty"`
	Ceil   uint64 `json:"ceil,omitempty"`
}
