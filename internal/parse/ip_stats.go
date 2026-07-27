package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// ipLinkStats mirrors the fields of `ip -s -j link show` ifscope consumes.
type ipLinkStats struct {
	IfIndex int    `json:"ifindex"`
	IfName  string `json:"ifname"`
	Stats64 *struct {
		Rx struct {
			Bytes   uint64 `json:"bytes"`
			Packets uint64 `json:"packets"`
			Errors  uint64 `json:"errors"`
			Dropped uint64 `json:"dropped"`
		} `json:"rx"`
		Tx struct {
			Bytes      uint64 `json:"bytes"`
			Packets    uint64 `json:"packets"`
			Errors     uint64 `json:"errors"`
			Dropped    uint64 `json:"dropped"`
			Collisions uint64 `json:"collisions"`
		} `json:"tx"`
	} `json:"stats64"`
}

// IPLinkStats parses `ip -s -j link show` into per-interface counters.
// Interfaces without a stats64 block are skipped.
func IPLinkStats(data []byte) ([]model.InterfaceStats, error) {
	var raw []ipLinkStats
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse ip -s link json: %w", err)
	}

	out := make([]model.InterfaceStats, 0, len(raw))
	for _, l := range raw {
		if l.Stats64 == nil {
			continue
		}
		out = append(out, model.InterfaceStats{
			ID:         l.IfIndex,
			Name:       l.IfName,
			RxBytes:    l.Stats64.Rx.Bytes,
			RxPackets:  l.Stats64.Rx.Packets,
			RxErrors:   l.Stats64.Rx.Errors,
			RxDropped:  l.Stats64.Rx.Dropped,
			TxBytes:    l.Stats64.Tx.Bytes,
			TxPackets:  l.Stats64.Tx.Packets,
			TxErrors:   l.Stats64.Tx.Errors,
			TxDropped:  l.Stats64.Tx.Dropped,
			Collisions: l.Stats64.Tx.Collisions,
		})
	}
	return out, nil
}
