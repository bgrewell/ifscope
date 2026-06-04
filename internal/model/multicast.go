package model

// MulticastGroup is an IP multicast group an interface has joined.
type MulticastGroup struct {
	Interface string `json:"interface"`
	Family    string `json:"family"`
	Address   string `json:"address"`
}
