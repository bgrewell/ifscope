package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// lldpValue is lldpcli's ubiquitous {"value": "..."} wrapper.
type lldpValue struct {
	Value string `json:"value"`
}

// lldpID is an id/{type,value} pair (chassis id, port id).
type lldpID struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// lldpJSON0 mirrors `lldpcli show neighbors -f json0`. The json0 format wraps
// every field in arrays, which (unlike -f json) is unambiguous regardless of
// how many interfaces/neighbors are present.
type lldpJSON0 struct {
	LLDP []struct {
		Interface []struct {
			Name    string `json:"name"`
			Chassis []struct {
				ID         []lldpID    `json:"id"`
				Name       []lldpValue `json:"name"`
				Descr      []lldpValue `json:"descr"`
				MgmtIP     []lldpValue `json:"mgmt-ip"`
				Capability []struct {
					Type    string `json:"type"`
					Enabled bool   `json:"enabled"`
				} `json:"capability"`
			} `json:"chassis"`
			Port []struct {
				ID    []lldpID    `json:"id"`
				Descr []lldpValue `json:"descr"`
				TTL   []lldpValue `json:"ttl"`
			} `json:"port"`
		} `json:"interface"`
	} `json:"lldp"`
}

// LLDPNeighbors parses `lldpcli show neighbors -f json0` into neighbor entries,
// one per local interface that has a neighbor.
func LLDPNeighbors(data []byte) ([]model.LLDPNeighbor, error) {
	var raw lldpJSON0
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse lldp json0: %w", err)
	}
	if len(raw.LLDP) == 0 {
		return nil, nil
	}

	var out []model.LLDPNeighbor
	for _, iface := range raw.LLDP[0].Interface {
		n := model.LLDPNeighbor{LocalPort: iface.Name}
		if len(iface.Chassis) > 0 {
			ch := iface.Chassis[0]
			n.Chassis = firstValue(ch.Name)
			if len(ch.ID) > 0 {
				n.ChassisID = ch.ID[0].Value
			}
			n.ChassisDescr = firstValue(ch.Descr)
			for _, ip := range ch.MgmtIP {
				n.MgmtIPs = append(n.MgmtIPs, ip.Value)
			}
			for _, c := range ch.Capability {
				if c.Enabled {
					n.Capabilities = append(n.Capabilities, c.Type)
				}
			}
		}
		if len(iface.Port) > 0 {
			p := iface.Port[0]
			if len(p.ID) > 0 {
				n.PortID = p.ID[0].Value
			}
			n.PortDescr = firstValue(p.Descr)
			n.TTL = firstValue(p.TTL)
		}
		out = append(out, n)
	}
	return out, nil
}

// firstValue returns the first wrapped value, or "".
func firstValue(vs []lldpValue) string {
	if len(vs) > 0 {
		return vs[0].Value
	}
	return ""
}
