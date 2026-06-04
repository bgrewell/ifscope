package model

// PTP describes an interface's hardware timestamping / PTP-clock capabilities
// (from ethtool -T).
type PTP struct {
	Name     string `json:"name"`
	PHCIndex *int   `json:"phc_index,omitempty"`
	HWTx     bool   `json:"hw_tx"`
	HWRx     bool   `json:"hw_rx"`
	SWTx     bool   `json:"sw_tx"`
	SWRx     bool   `json:"sw_rx"`
}
