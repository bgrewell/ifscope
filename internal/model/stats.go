package model

// InterfaceStats holds per-interface traffic and error counters.
type InterfaceStats struct {
	ID           int                 `json:"id,omitempty"`
	Name         string              `json:"name"`
	RxBytes      uint64              `json:"rx_bytes"`
	RxPackets    uint64              `json:"rx_packets"`
	RxErrors     uint64              `json:"rx_errors"`
	RxDropped    uint64              `json:"rx_dropped"`
	TxBytes      uint64              `json:"tx_bytes"`
	TxPackets    uint64              `json:"tx_packets"`
	TxErrors     uint64              `json:"tx_errors"`
	TxDropped    uint64              `json:"tx_dropped"`
	Collisions   uint64              `json:"collisions,omitempty"`
	LinkSpeedBps uint64              `json:"link_speed_bps,omitempty"`
	Rates        *InterfaceStatsRate `json:"rates,omitempty"`
	RateStatus   string              `json:"rate_status,omitempty"`
}

// InterfaceStatsRate holds per-second rates derived from two counter samples.
// Raw counters remain on InterfaceStats so machine consumers can use either.
type InterfaceStatsRate struct {
	SampleSeconds float64 `json:"sample_seconds"`
	RxBitsPerSec  float64 `json:"rx_bits_per_sec"`
	RxPacketsSec  float64 `json:"rx_packets_per_sec"`
	RxErrorsSec   float64 `json:"rx_errors_per_sec"`
	RxDroppedSec  float64 `json:"rx_dropped_per_sec"`
	TxBitsPerSec  float64 `json:"tx_bits_per_sec"`
	TxPacketsSec  float64 `json:"tx_packets_per_sec"`
	TxErrorsSec   float64 `json:"tx_errors_per_sec"`
	TxDroppedSec  float64 `json:"tx_dropped_per_sec"`
	RxUtilization float64 `json:"rx_utilization_percent,omitempty"`
	TxUtilization float64 `json:"tx_utilization_percent,omitempty"`
}
