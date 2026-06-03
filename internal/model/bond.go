package model

// Bond is a Linux bonding (link aggregation) master and its enslaved members.
type Bond struct {
	Name        string       `json:"name"`
	Mode        string       `json:"mode,omitempty"`
	State       string       `json:"state"`
	ActiveSlave string       `json:"active_slave,omitempty"`
	Members     []BondMember `json:"members,omitempty"`
}

// BondMember is one interface enslaved to a bond, with its link state.
type BondMember struct {
	Name  string `json:"name"`
	State string `json:"state"`
}
