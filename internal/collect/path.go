package collect

import (
	"context"
	"fmt"
	"net"
	"net/netip"
	"sort"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// ResolveFunc resolves a hostname or address into destination addresses.
type ResolveFunc func(context.Context, string) ([]netip.Addr, error)

// PathOptions controls a passive kernel route lookup.
type PathOptions struct {
	Family   string
	Source   string
	OutIface string
	Protocol string
	Port     int
}

// Path collects route decisions and correlates them with local network state.
type Path struct {
	Runner  run.Runner
	Resolve ResolveFunc
}

// NewPath returns a path collector using the system resolver.
func NewPath(r run.Runner) *Path {
	return &Path{Runner: r, Resolve: resolveAddresses}
}

// Collect resolves destination and asks the kernel for the route to every
// candidate. It does not send traffic.
func (c *Path) Collect(ctx context.Context, destination string, opts PathOptions) (model.Path, []model.Warning) {
	result := model.Path{Destination: destination, Port: opts.Port, Protocol: opts.Protocol}
	addrs, err := c.Resolve(ctx, destination)
	if err != nil {
		return result, []model.Warning{{Source: "dns", Message: fmt.Sprintf("resolve %q: %v", destination, err), Fatal: true}}
	}
	addrs = filterAddresses(addrs, opts.Family)
	if len(addrs) == 0 {
		return result, []model.Warning{{Source: "dns", Message: fmt.Sprintf("no %s addresses found for %q", familyLabel(opts.Family), destination), Fatal: true}}
	}

	ifaces, iw := NewInterfaces(c.Runner).Collect(ctx)
	routes, tw := NewRoutes(c.Runner).Collect(ctx)
	rules, rw := NewRules(c.Runner).Collect(ctx)
	neighbors, nw := NewNeighbors(c.Runner).Collect(ctx)
	warnings := append(append(append(iw, tw...), rw...), nw...)
	if containsIPv6(addrs) {
		v6rules, vw := c.ipv6Rules(ctx)
		rules = append(rules, v6rules...)
		warnings = append(warnings, vw...)
	}

	for _, addr := range addrs {
		candidate := model.PathCandidate{Address: addr.String(), Family: addrFamily(addr)}
		route, lookupErr := c.lookup(ctx, addr, opts)
		if lookupErr != nil {
			candidate.Error = lookupErr.Error()
			result.Candidates = append(result.Candidates, candidate)
			continue
		}
		if matched, matchErr := c.fibMatch(ctx, addr, opts); matchErr == nil {
			route = correlateRoute(route, matched, routes)
		}
		candidate.Route = &route
		source := opts.Source
		if source == "" {
			source = route.Src
		}
		if rule := selectedRule(rules, route, addr, source, opts.OutIface); rule != nil {
			candidate.Rule = rule
		}
		candidate.Topology, candidate.MTU = topology(ifaces, route.Dev)
		next := route.Gateway
		if next == "" {
			next = addr.String()
		}
		for i := range neighbors {
			if neighbors[i].Dst == next && (route.Dev == "" || neighbors[i].Dev == route.Dev) {
				n := neighbors[i]
				candidate.Neighbor = &n
				break
			}
		}
		result.Candidates = append(result.Candidates, candidate)
	}
	return result, warnings
}

func (c *Path) ipv6Rules(ctx context.Context) ([]model.Rule, []model.Warning) {
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, "-6", "-json", "rule")
	if err != nil {
		return nil, []model.Warning{{
			Source: "ip", Message: fmt.Sprintf("IPv6 policy rules unavailable: %v: %s", err, strings.TrimSpace(string(stderr))),
		}}
	}
	rules, err := parse.IPRules(stdout)
	if err != nil {
		return nil, []model.Warning{{Source: "ip", Message: err.Error()}}
	}
	for i := range rules {
		rules[i].Family = "inet6"
	}
	return rules, nil
}

func (c *Path) lookup(ctx context.Context, addr netip.Addr, opts PathOptions) (model.Route, error) {
	args := routeGetArgs(addr, opts)
	stdout, stderr, err := c.Runner.Run(ctx, ipCmd, args...)
	if err != nil {
		if run.IsNotFound(err) {
			return model.Route{}, fmt.Errorf("ip command not found")
		}
		return model.Route{}, fmt.Errorf("ip route get failed: %v: %s", err, strings.TrimSpace(string(stderr)))
	}
	routes, err := parse.IPRoutes(stdout)
	if err != nil {
		return model.Route{}, err
	}
	if len(routes) == 0 {
		return model.Route{}, fmt.Errorf("kernel returned no route")
	}
	return routes[0], nil
}

func (c *Path) fibMatch(ctx context.Context, addr netip.Addr, opts PathOptions) (model.Route, error) {
	args := append(routeGetArgs(addr, opts), "fibmatch")
	stdout, _, err := c.Runner.Run(ctx, ipCmd, args...)
	if err != nil {
		return model.Route{}, err
	}
	routes, err := parse.IPRoutes(stdout)
	if err != nil || len(routes) == 0 {
		return model.Route{}, fmt.Errorf("kernel returned no matching route")
	}
	return routes[0], nil
}

func routeGetArgs(addr netip.Addr, opts PathOptions) []string {
	args := []string{"-json"}
	if addr.Is4() {
		args = append(args, "-4")
	} else {
		args = append(args, "-6")
	}
	args = append(args, "route", "get", addr.String())
	if opts.Source != "" {
		args = append(args, "from", opts.Source)
	}
	if opts.OutIface != "" {
		args = append(args, "oif", opts.OutIface)
	}
	if opts.Protocol != "" {
		args = append(args, "ipproto", opts.Protocol)
	}
	if opts.Port != 0 {
		args = append(args, "dport", strconv.Itoa(opts.Port))
	}
	return args
}

// correlateRoute combines the selected source/gateway from route-get with the
// matched FIB prefix and table from the full route inventory.
func correlateRoute(selected, fib model.Route, inventory []model.Route) model.Route {
	selected.Dst = fib.Dst
	selected.Protocol = fib.Protocol
	selected.Metric = fib.Metric
	selected.Scope = fib.Scope
	if fib.Table != "" {
		selected.Table = fib.Table
	}
	for _, route := range inventory {
		if route.Family != fib.Family || route.Dst != fib.Dst {
			continue
		}
		if fib.Dev != "" && route.Dev != fib.Dev {
			continue
		}
		if fib.Gateway != "" && route.Gateway != fib.Gateway {
			continue
		}
		selected.Table = route.Table
		break
	}
	if selected.Table == "" {
		selected.Table = "main"
	}
	return selected
}

func resolveAddresses(ctx context.Context, destination string) ([]netip.Addr, error) {
	if addr, err := netip.ParseAddr(destination); err == nil {
		return []netip.Addr{addr.Unmap()}, nil
	}
	addrs, err := net.DefaultResolver.LookupNetIP(ctx, "ip", destination)
	if err != nil {
		return nil, err
	}
	seen := map[netip.Addr]bool{}
	out := make([]netip.Addr, 0, len(addrs))
	for _, addr := range addrs {
		addr = addr.Unmap()
		if addr.IsValid() && !seen[addr] {
			seen[addr] = true
			out = append(out, addr)
		}
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Is4() != out[j].Is4() {
			return out[i].Is4()
		}
		return out[i].Less(out[j])
	})
	return out, nil
}

func filterAddresses(addrs []netip.Addr, family string) []netip.Addr {
	if family == "" || family == "any" {
		return addrs
	}
	out := make([]netip.Addr, 0, len(addrs))
	for _, addr := range addrs {
		if (family == "4" && addr.Is4()) || (family == "6" && addr.Is6()) {
			out = append(out, addr)
		}
	}
	return out
}

func containsIPv6(addrs []netip.Addr) bool {
	for _, addr := range addrs {
		if addr.Is6() {
			return true
		}
	}
	return false
}

func addrFamily(addr netip.Addr) string {
	if addr.Is6() {
		return "inet6"
	}
	return "inet"
}

func familyLabel(family string) string {
	switch family {
	case "4":
		return "IPv4"
	case "6":
		return "IPv6"
	default:
		return "IP"
	}
}

func topology(ifaces []model.Interface, dev string) ([]model.PathHop, int) {
	byName := make(map[string]model.Interface, len(ifaces))
	for _, iface := range ifaces {
		byName[iface.Name] = iface
	}
	var hops []model.PathHop
	mtu := 0
	seen := map[string]bool{}
	for dev != "" && !seen[dev] {
		seen[dev] = true
		iface, ok := byName[dev]
		if !ok {
			break
		}
		hops = append(hops, model.PathHop{Name: iface.Name, Type: iface.Type, MTU: iface.MTU})
		if iface.MTU > 0 && (mtu == 0 || iface.MTU < mtu) {
			mtu = iface.MTU
		}
		dev = iface.LinkParent
	}
	return hops, mtu
}

func selectedRule(rules []model.Rule, route model.Route, dst netip.Addr, source, oif string) *model.Rule {
	for i := range rules {
		r := &rules[i]
		if r.Family != "" && r.Family != addrFamily(dst) {
			continue
		}
		if r.Table != "" && r.Table != route.Table {
			continue
		}
		if r.OIf != "" && r.OIf != oif && r.OIf != route.Dev {
			continue
		}
		if !matchesAddress(r.To, dst) || !matchesStringAddress(r.From, source) {
			continue
		}
		copy := *r
		return &copy
	}
	return nil
}

func matchesStringAddress(prefix, value string) bool {
	if prefix == "" || prefix == "all" {
		return true
	}
	addr, err := netip.ParseAddr(value)
	return err == nil && matchesAddress(prefix, addr)
}

func matchesAddress(value string, addr netip.Addr) bool {
	if value == "" || value == "all" {
		return true
	}
	if prefix, err := netip.ParsePrefix(value); err == nil {
		return prefix.Contains(addr)
	}
	want, err := netip.ParseAddr(value)
	return err == nil && want == addr
}
