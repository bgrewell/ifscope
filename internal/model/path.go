package model

// Path describes how the local host would send traffic to a destination.
type Path struct {
	Destination string          `json:"destination"`
	Port        int             `json:"port,omitempty"`
	Protocol    string          `json:"protocol,omitempty"`
	Candidates  []PathCandidate `json:"candidates"`
}

// PathCandidate is the kernel-selected path for one resolved destination
// address. Error is set when lookup failed for this candidate.
type PathCandidate struct {
	Address  string    `json:"address"`
	Family   string    `json:"family"`
	Route    *Route    `json:"route,omitempty"`
	Rule     *Rule     `json:"rule,omitempty"`
	Neighbor *Neighbor `json:"neighbor,omitempty"`
	Topology []PathHop `json:"topology,omitempty"`
	MTU      int       `json:"mtu,omitempty"`
	Error    string    `json:"error,omitempty"`
}

// PathHop is one local interface in the egress-to-parent topology chain.
type PathHop struct {
	Name string        `json:"name"`
	Type InterfaceType `json:"type"`
	MTU  int           `json:"mtu,omitempty"`
}
