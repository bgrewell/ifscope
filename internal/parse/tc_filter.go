package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// tcFilter mirrors `tc -json filter show dev <d>`. The command emits several
// entries per filter (a header, a hash-table entry, then the rule); only the
// rule carries options.flowid, which is the entry ifscope keeps.
type tcFilter struct {
	Parent   string `json:"parent"`
	Protocol string `json:"protocol"`
	Pref     int    `json:"pref"`
	Kind     string `json:"kind"`
	Options  *struct {
		FlowID string `json:"flowid"`
	} `json:"options"`
}

// TCFilters parses `tc -json filter show dev <d>`, keeping only the rule entries
// that classify into a flow (class). The Dev field is set by the caller. An
// error is returned when the output is not JSON (older tc emits text).
func TCFilters(data []byte) ([]model.TCFilter, error) {
	var raw []tcFilter
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse tc filter json: %w", err)
	}
	var out []model.TCFilter
	for _, f := range raw {
		if f.Options == nil || f.Options.FlowID == "" {
			continue
		}
		out = append(out, model.TCFilter{
			Parent:   f.Parent,
			Protocol: f.Protocol,
			Pref:     f.Pref,
			Kind:     f.Kind,
			FlowID:   f.Options.FlowID,
		})
	}
	return out, nil
}
