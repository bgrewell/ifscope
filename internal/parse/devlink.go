package parse

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/bgrewell/ifscope/internal/model"
)

// devlinkPortShow mirrors `devlink -j port show`: an object keyed by port
// handle (e.g. "pci/0000:17:00.0/0").
type devlinkPortShow struct {
	Port map[string]struct {
		Type    string `json:"type"`
		Netdev  string `json:"netdev"`
		Flavour string `json:"flavour"`
		PfNum   *int   `json:"pfnum"`
		VfNum   *int   `json:"vfnum"`
		Lanes   int    `json:"lanes"`
	} `json:"port"`
}

// DevlinkPorts parses `devlink -j port show` into ports sorted by handle.
func DevlinkPorts(data []byte) ([]model.DevlinkPort, error) {
	var raw devlinkPortShow
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse devlink port json: %w", err)
	}

	handles := make([]string, 0, len(raw.Port))
	for h := range raw.Port {
		handles = append(handles, h)
	}
	sort.Strings(handles)

	out := make([]model.DevlinkPort, 0, len(handles))
	for _, h := range handles {
		p := raw.Port[h]
		out = append(out, model.DevlinkPort{
			Handle:  h,
			Type:    p.Type,
			Flavour: p.Flavour,
			Netdev:  p.Netdev,
			PfNum:   p.PfNum,
			VfNum:   p.VfNum,
			Lanes:   p.Lanes,
		})
	}
	return out, nil
}
