package model

// TCFilter is a traffic-control filter that classifies traffic into a class
// (FlowID), e.g. a u32 or flower rule under an htb/hfsc qdisc.
type TCFilter struct {
	Dev      string `json:"dev"`
	Parent   string `json:"parent,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Pref     int    `json:"pref,omitempty"`
	Kind     string `json:"kind"`
	FlowID   string `json:"flowid"`
}
