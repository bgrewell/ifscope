package model

// DNS is the resolver state for a single link (or the global scope, when Link
// is empty/"global"), as reported by resolvectl.
type DNS struct {
	Link          string   `json:"link"`
	CurrentServer string   `json:"current_server,omitempty"`
	Servers       []string `json:"servers,omitempty"`
	Domains       []string `json:"domains,omitempty"`
	DefaultRoute  *bool    `json:"default_route,omitempty"`
	LLMNR         string   `json:"llmnr,omitempty"`
	MDNS          string   `json:"mdns,omitempty"`
	DNSSEC        string   `json:"dnssec,omitempty"`
}
