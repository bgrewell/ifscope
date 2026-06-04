package model

// Offloads holds an interface's NIC offload feature states (from ethtool -k).
// Features maps each top-level feature name to its state ("on", "off",
// "off [fixed]"). The human table surfaces the common ones as columns.
type Offloads struct {
	Name     string            `json:"name"`
	Features map[string]string `json:"features,omitempty"`
}
