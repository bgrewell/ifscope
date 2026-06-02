package model

// Address is a single IP address assigned to an interface.
//
// Family is "inet" for IPv4 and "inet6" for IPv6. The JSON schema always
// carries both families; human table rendering filters to IPv4 by default.
type Address struct {
	Family    string `json:"family"`
	Local     string `json:"local"`
	PrefixLen int    `json:"prefixlen"`
	Scope     string `json:"scope,omitempty"`
	Dynamic   bool   `json:"dynamic,omitempty"`
	Secondary bool   `json:"secondary,omitempty"`
}

// IsIPv4 reports whether the address belongs to the IPv4 family.
func (a Address) IsIPv4() bool { return a.Family == "inet" }

// IsIPv6 reports whether the address belongs to the IPv6 family.
func (a Address) IsIPv6() bool { return a.Family == "inet6" }
