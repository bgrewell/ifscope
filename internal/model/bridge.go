package model

// Bridge is a native Linux bridge and its member ports. (Open vSwitch bridges
// are reported separately under OVS.)
type Bridge struct {
	Name    string         `json:"name"`
	State   string         `json:"state"`
	STP     bool           `json:"stp"`
	Members []BridgeMember `json:"members,omitempty"`
}

// BridgeMember is an interface enslaved to a Linux bridge, with its link state.
type BridgeMember struct {
	Name  string `json:"name"`
	State string `json:"state"`
}
