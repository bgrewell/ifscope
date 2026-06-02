package parse

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// ipRoute mirrors one element of `ip -detail -json route`.
type ipRoute struct {
	Type     string   `json:"type"`
	Dst      string   `json:"dst"`
	Gateway  string   `json:"gateway"`
	Dev      string   `json:"dev"`
	Protocol string   `json:"protocol"`
	Scope    string   `json:"scope"`
	PrefSrc  string   `json:"prefsrc"`
	Metric   int      `json:"metric"`
	Table    string   `json:"table"`
	Flags    []string `json:"flags"`
}

// IPRoutes parses `ip -detail -json route` into routes. Address family is
// inferred from the addresses present, since iproute2 does not emit it here.
func IPRoutes(data []byte) ([]model.Route, error) {
	var raw []ipRoute
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse ip route json: %w", err)
	}

	out := make([]model.Route, 0, len(raw))
	for _, r := range raw {
		out = append(out, model.Route{
			Dst:      r.Dst,
			Gateway:  r.Gateway,
			Dev:      r.Dev,
			Protocol: r.Protocol,
			Metric:   r.Metric,
			Table:    r.Table,
			Scope:    r.Scope,
			Src:      r.PrefSrc,
			Family:   routeFamily(r),
		})
	}
	return out, nil
}

// routeFamily infers inet6 when any address in the route is IPv6.
func routeFamily(r ipRoute) string {
	if strings.Contains(r.Dst, ":") ||
		strings.Contains(r.Gateway, ":") ||
		strings.Contains(r.PrefSrc, ":") {
		return "inet6"
	}
	return "inet"
}
