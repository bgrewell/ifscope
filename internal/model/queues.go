package model

// Count is a current/maximum pair (e.g. configured vs supported queues).
type Count struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

// Queues holds an interface's channel/ring/coalesce/RSS/RPS-XPS configuration
// from ethtool (-l, -g, -c, -x) and sysfs queue masks.
type Queues struct {
	Name       string `json:"name"`
	Combined   Count  `json:"combined"`
	RxChannels Count  `json:"rx_channels"`
	TxChannels Count  `json:"tx_channels"`
	RxRing     Count  `json:"rx_ring"`
	TxRing     Count  `json:"tx_ring"`
	RxUsecs    int    `json:"rx_usecs,omitempty"`
	TxUsecs    int    `json:"tx_usecs,omitempty"`
	AdaptiveRx bool   `json:"adaptive_rx,omitempty"`
	AdaptiveTx bool   `json:"adaptive_tx,omitempty"`
	RSSRings   int    `json:"rss_rings,omitempty"`
	// RPSQueues/XPSQueues count the rx/tx queues with a non-zero steering mask.
	RPSQueues int `json:"rps_queues,omitempty"`
	XPSQueues int `json:"xps_queues,omitempty"`
}
