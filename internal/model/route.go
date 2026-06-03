package model

// Route is a single entry from the kernel routing tables.
//
// Dst is "default" for the default route. Family is "inet" or "inet6".
type Route struct {
	Dst      string `json:"dst"`
	Gateway  string `json:"gateway,omitempty"`
	Dev      string `json:"dev,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Metric   int    `json:"metric,omitempty"`
	Table    string `json:"table,omitempty"`
	Scope    string `json:"scope,omitempty"`
	Src      string `json:"src,omitempty"`
	Family   string `json:"family,omitempty"`
	// NextHops carries the per-next-hop details of a multipath (ECMP) route.
	// For such routes the top-level Gateway/Dev are empty.
	NextHops []RouteNextHop `json:"nexthops,omitempty"`
}

// RouteNextHop is one path of a multipath route.
type RouteNextHop struct {
	Gateway string `json:"gateway,omitempty"`
	Dev     string `json:"dev,omitempty"`
	Weight  int    `json:"weight,omitempty"`
}
