package model

// Count is a current/maximum pair (e.g. configured vs supported queues).
type Count struct {
	Current int `json:"current"`
	Max     int `json:"max"`
}

// Queues holds an interface's channel (queue) and ring-buffer configuration
// from ethtool -l / -g.
type Queues struct {
	Name       string `json:"name"`
	Combined   Count  `json:"combined"`
	RxChannels Count  `json:"rx_channels"`
	TxChannels Count  `json:"tx_channels"`
	RxRing     Count  `json:"rx_ring"`
	TxRing     Count  `json:"tx_ring"`
}
