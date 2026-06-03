package parse

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// ipRule mirrors one element of `ip -json rule`. The kernel reports the match
// source/destination as src/dst; ifscope exposes them as from/to.
type ipRule struct {
	Priority int    `json:"priority"`
	Src      string `json:"src"`
	SrcLen   int    `json:"srclen"`
	Dst      string `json:"dst"`
	DstLen   int    `json:"dstlen"`
	Table    string `json:"table"`
	FWMark   string `json:"fwmark"`
	IIf      string `json:"iif"`
	OIf      string `json:"oif"`
	Protocol string `json:"protocol"`
}

// IPRules parses `ip -json rule` into policy rules.
func IPRules(data []byte) ([]model.Rule, error) {
	var raw []ipRule
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse ip rule json: %w", err)
	}

	out := make([]model.Rule, 0, len(raw))
	for _, r := range raw {
		out = append(out, model.Rule{
			Priority: r.Priority,
			From:     prefix(r.Src, r.SrcLen),
			To:       prefix(r.Dst, r.DstLen),
			Table:    r.Table,
			FWMark:   r.FWMark,
			IIf:      r.IIf,
			OIf:      r.OIf,
			Family:   ruleFamily(r),
		})
	}
	return out, nil
}

// prefix joins an address with its prefix length when the length narrows it.
// "all" and bare host addresses are returned unchanged.
func prefix(addr string, length int) string {
	if addr == "" || addr == "all" {
		return addr
	}
	if length == 0 || strings.Contains(addr, "/") {
		return addr
	}
	return fmt.Sprintf("%s/%d", addr, length)
}

// ruleFamily infers inet6 from any IPv6 address in the rule.
func ruleFamily(r ipRule) string {
	if strings.Contains(r.Src, ":") || strings.Contains(r.Dst, ":") {
		return "inet6"
	}
	return "inet"
}
