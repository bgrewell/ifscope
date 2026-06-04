package model

// Socket is a listening socket (TCP/UDP). Process is populated only when the
// command runs with sufficient privilege.
type Socket struct {
	Proto     string `json:"proto"`
	State     string `json:"state,omitempty"`
	LocalAddr string `json:"local_addr"`
	LocalPort string `json:"local_port"`
	Process   string `json:"process,omitempty"`
}
