package parse

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// ipNeigh mirrors one element of `ip -json neigh show`. State is an array
// (usually one entry, e.g. ["REACHABLE"]).
type ipNeigh struct {
	Dst    string   `json:"dst"`
	Dev    string   `json:"dev"`
	LLAddr string   `json:"lladdr"`
	State  []string `json:"state"`
	Router bool     `json:"router"`
}

// IPNeighbors parses `ip -json neigh show` into neighbor entries. Family is
// inferred from the destination address.
func IPNeighbors(data []byte) ([]model.Neighbor, error) {
	var raw []ipNeigh
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse ip neigh json: %w", err)
	}

	out := make([]model.Neighbor, 0, len(raw))
	for _, n := range raw {
		family := "inet"
		if strings.Contains(n.Dst, ":") {
			family = "inet6"
		}
		out = append(out, model.Neighbor{
			Dst:    n.Dst,
			Dev:    n.Dev,
			LLAddr: n.LLAddr,
			State:  strings.Join(n.State, ","),
			Router: n.Router,
			Family: family,
		})
	}
	return out, nil
}
