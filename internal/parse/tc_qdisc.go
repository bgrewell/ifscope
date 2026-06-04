package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// tcQdisc mirrors one element of `tc -json qdisc show`.
type tcQdisc struct {
	Kind   string `json:"kind"`
	Handle string `json:"handle"`
	Dev    string `json:"dev"`
	Parent string `json:"parent"`
	Root   bool   `json:"root"`
}

// Qdiscs parses `tc -json qdisc show` into all qdiscs (root and child).
func Qdiscs(data []byte) ([]model.Qdisc, error) {
	var raw []tcQdisc
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse tc qdisc json: %w", err)
	}
	out := make([]model.Qdisc, 0, len(raw))
	for _, q := range raw {
		out = append(out, model.Qdisc{
			Dev:    q.Dev,
			Kind:   q.Kind,
			Handle: q.Handle,
			Parent: q.Parent,
			Root:   q.Root,
		})
	}
	return out, nil
}
