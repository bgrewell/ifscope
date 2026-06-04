package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// bridgeMDBGroup mirrors `bridge -json mdb show`: an array of objects each
// carrying an "mdb" list of entries.
type bridgeMDBGroup struct {
	MDB []struct {
		Dev   string `json:"dev"`
		Port  string `json:"port"`
		Grp   string `json:"grp"`
		Vid   int    `json:"vid"`
		State string `json:"state"`
	} `json:"mdb"`
}

// BridgeMDB parses `bridge -json mdb show` into multicast-database entries.
func BridgeMDB(data []byte) ([]model.MDBEntry, error) {
	var raw []bridgeMDBGroup
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse bridge mdb json: %w", err)
	}
	var out []model.MDBEntry
	for _, g := range raw {
		for _, e := range g.MDB {
			out = append(out, model.MDBEntry{
				Bridge: e.Dev,
				Port:   e.Port,
				Group:  e.Grp,
				VLAN:   e.Vid,
				State:  e.State,
			})
		}
	}
	return out, nil
}
