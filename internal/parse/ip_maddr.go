package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// ipMaddr mirrors one element of `ip -json maddr show`.
type ipMaddr struct {
	IfName string `json:"ifname"`
	Maddr  []struct {
		Family  string `json:"family"`
		Address string `json:"address"`
	} `json:"maddr"`
}

// MulticastGroups parses `ip -json maddr show` into IP multicast group
// memberships. Pure link-layer (MAC) entries, which carry no "family", are
// skipped in favor of the inet/inet6 groups.
func MulticastGroups(data []byte) ([]model.MulticastGroup, error) {
	var raw []ipMaddr
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse ip maddr json: %w", err)
	}
	var out []model.MulticastGroup
	for _, e := range raw {
		for _, m := range e.Maddr {
			if m.Family == "" || m.Address == "" {
				continue
			}
			out = append(out, model.MulticastGroup{
				Interface: e.IfName,
				Family:    m.Family,
				Address:   m.Address,
			})
		}
	}
	return out, nil
}
