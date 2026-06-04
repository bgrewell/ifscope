package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// tcClass mirrors one element of `tc -json class show dev <d>`. rate/ceil are in
// bytes per second.
type tcClass struct {
	Class  string `json:"class"`
	Handle string `json:"handle"`
	Parent string `json:"parent"`
	Root   bool   `json:"root"`
	Rate   uint64 `json:"rate"`
	Ceil   uint64 `json:"ceil"`
}

// TCClasses parses `tc -json class show dev <d>`. rate/ceil are converted from
// bytes/s to bits/s. The Dev field is left for the caller to set. An error is
// returned when the output is not JSON (older tc emits text for `class show`).
func TCClasses(data []byte) ([]model.TCClass, error) {
	var raw []tcClass
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse tc class json: %w", err)
	}
	out := make([]model.TCClass, 0, len(raw))
	for _, c := range raw {
		out = append(out, model.TCClass{
			Kind:   c.Class,
			Handle: c.Handle,
			Parent: c.Parent,
			Root:   c.Root,
			Rate:   c.Rate * 8,
			Ceil:   c.Ceil * 8,
		})
	}
	return out, nil
}
