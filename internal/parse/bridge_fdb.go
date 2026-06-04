package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// bridgeFDB mirrors one element of `bridge -json fdb show`.
type bridgeFDB struct {
	MAC    string   `json:"mac"`
	IfName string   `json:"ifname"`
	VLAN   *int     `json:"vlan"`
	Master string   `json:"master"`
	State  string   `json:"state"`
	Flags  []string `json:"flags"`
}

// BridgeFDB parses `bridge -json fdb show` into forwarding-database entries.
func BridgeFDB(data []byte) ([]model.FDBEntry, error) {
	var raw []bridgeFDB
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse bridge fdb json: %w", err)
	}
	out := make([]model.FDBEntry, 0, len(raw))
	for _, e := range raw {
		out = append(out, model.FDBEntry{
			MAC:    e.MAC,
			Dev:    e.IfName,
			VLAN:   e.VLAN,
			Master: e.Master,
			State:  e.State,
			Flags:  e.Flags,
		})
	}
	return out, nil
}
