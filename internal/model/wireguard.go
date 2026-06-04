package model

// WireGuardDevice is a WireGuard interface and its peers. Private and
// pre-shared keys are intentionally never collected.
type WireGuardDevice struct {
	Interface  string   `json:"interface"`
	PublicKey  string   `json:"public_key,omitempty"`
	ListenPort int      `json:"listen_port,omitempty"`
	Peers      []WGPeer `json:"peers,omitempty"`
}

// WGPeer is a WireGuard peer.
type WGPeer struct {
	PublicKey       string   `json:"public_key"`
	Endpoint        string   `json:"endpoint,omitempty"`
	AllowedIPs      []string `json:"allowed_ips,omitempty"`
	LatestHandshake int64    `json:"latest_handshake,omitempty"`
	RxBytes         uint64   `json:"rx_bytes,omitempty"`
	TxBytes         uint64   `json:"tx_bytes,omitempty"`
	Keepalive       string   `json:"keepalive,omitempty"`
}
