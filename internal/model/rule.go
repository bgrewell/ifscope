package model

// Rule is a routing policy rule (`ip rule`) that selects which routing table is
// consulted for matching traffic — the basis of source-based / policy routing.
type Rule struct {
	Priority int    `json:"priority"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Table    string `json:"table,omitempty"`
	FWMark   string `json:"fwmark,omitempty"`
	IIf      string `json:"iif,omitempty"`
	OIf      string `json:"oif,omitempty"`
	Family   string `json:"family,omitempty"`
}
