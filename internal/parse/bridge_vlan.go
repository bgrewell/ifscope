package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// bridgeVLANPort mirrors one element of `bridge -json vlan show`.
type bridgeVLANPort struct {
	IfName string `json:"ifname"`
	Vlans  []struct {
		Vlan    int      `json:"vlan"`
		VlanEnd int      `json:"vlanEnd"`
		Flags   []string `json:"flags"`
	} `json:"vlans"`
}

// BridgeVLANs parses `bridge -json vlan show` into per-port VLAN entries.
func BridgeVLANs(data []byte) ([]model.BridgeVLAN, error) {
	var raw []bridgeVLANPort
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse bridge vlan json: %w", err)
	}
	var out []model.BridgeVLAN
	for _, p := range raw {
		for _, v := range p.Vlans {
			out = append(out, model.BridgeVLAN{
				Port:    p.IfName,
				VLAN:    v.Vlan,
				VLANEnd: v.VlanEnd,
				Flags:   v.Flags,
			})
		}
	}
	return out, nil
}
