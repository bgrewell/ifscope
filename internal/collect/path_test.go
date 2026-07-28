package collect

import (
	"context"
	"net/netip"
	"reflect"
	"strings"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
)

func TestPathCollectCorrelatesRouteRuleNeighborAndTopology(t *testing.T) {
	fake := run.NewFake().
		Set(`[{"dst":"203.0.113.10","gateway":"192.0.2.1","dev":"vlan10","prefsrc":"192.0.2.10","table":"100"}]`,
			"ip", "-json", "-4", "route", "get", "203.0.113.10", "from", "192.0.2.10", "ipproto", "tcp", "dport", "443").
		Set(`[{"dst":"203.0.113.0/24","gateway":"192.0.2.1","dev":"vlan10","protocol":"static","table":"100"}]`,
			"ip", "-json", "-4", "route", "get", "203.0.113.10", "from", "192.0.2.10", "ipproto", "tcp", "dport", "443", "fibmatch").
		Set(`[
			{"ifindex":2,"ifname":"eth0","mtu":9000,"operstate":"UP","parentbus":"pci","parentdev":"0000:01:00.0"},
			{"ifindex":3,"ifname":"vlan10","mtu":1500,"operstate":"UP","link":"eth0","linkinfo":{"info_kind":"vlan","info_data":{"id":10}}}
		]`, "ip", "-detail", "-json", "address", "show").
		Set(`[{"dst":"203.0.113.0/24","gateway":"192.0.2.1","dev":"vlan10","protocol":"static","table":"100"}]`,
			"ip", "-detail", "-json", "route", "show", "table", "all").
		Set(`[{"priority":100,"src":"192.0.2.0","srclen":24,"table":"100"}]`, "ip", "-json", "rule").
		Set(`[{"dst":"192.0.2.1","dev":"vlan10","lladdr":"00:11:22:33:44:55","state":["REACHABLE"]}]`,
			"ip", "-json", "neigh", "show")

	collector := NewPath(fake)
	collector.Resolve = func(context.Context, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("203.0.113.10")}, nil
	}
	got, warnings := collector.Collect(context.Background(), "example.test", PathOptions{
		Source: "192.0.2.10", Protocol: "tcp", Port: 443,
	})

	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %+v", warnings)
	}
	if len(got.Candidates) != 1 {
		t.Fatalf("got %d candidates, want 1", len(got.Candidates))
	}
	candidate := got.Candidates[0]
	if candidate.Route == nil || candidate.Route.Dev != "vlan10" || candidate.Route.Table != "100" {
		t.Fatalf("unexpected route: %+v", candidate.Route)
	}
	if candidate.Rule == nil || candidate.Rule.Priority != 100 {
		t.Fatalf("unexpected rule: %+v", candidate.Rule)
	}
	if candidate.Neighbor == nil || candidate.Neighbor.State != "REACHABLE" {
		t.Fatalf("unexpected neighbor: %+v", candidate.Neighbor)
	}
	if candidate.MTU != 1500 {
		t.Errorf("MTU = %d, want 1500", candidate.MTU)
	}
	var names []string
	for _, hop := range candidate.Topology {
		names = append(names, hop.Name)
	}
	if !reflect.DeepEqual(names, []string{"vlan10", "eth0"}) {
		t.Errorf("topology = %v, want [vlan10 eth0]", names)
	}
}

func TestPathCollectRetainsPerCandidateLookupFailure(t *testing.T) {
	fake := run.NewFake().
		Set(`[]`, "ip", "-detail", "-json", "address", "show").
		Set(`[]`, "ip", "-detail", "-json", "route", "show", "table", "all").
		Set(`[]`, "ip", "-json", "rule").
		Set(`[]`, "ip", "-json", "neigh", "show")
	collector := NewPath(fake)
	collector.Resolve = func(context.Context, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("2001:db8::1")}, nil
	}

	got, _ := collector.Collect(context.Background(), "unreachable.test", PathOptions{Family: "6"})
	if len(got.Candidates) != 1 || !strings.Contains(got.Candidates[0].Error, "ip command not found") {
		t.Fatalf("unexpected candidates: %+v", got.Candidates)
	}
}

func TestFilterAddresses(t *testing.T) {
	addrs := []netip.Addr{
		netip.MustParseAddr("192.0.2.1"),
		netip.MustParseAddr("2001:db8::1"),
	}
	if got := filterAddresses(addrs, "4"); len(got) != 1 || !got[0].Is4() {
		t.Errorf("IPv4 filter = %v", got)
	}
	if got := filterAddresses(addrs, "6"); len(got) != 1 || !got[0].Is6() {
		t.Errorf("IPv6 filter = %v", got)
	}
}

func TestPathCollectUsesIPv6Rules(t *testing.T) {
	fake := run.NewFake().
		Set(`[]`, "ip", "-detail", "-json", "address", "show").
		Set(`[{"dst":"default","dev":"eth0","table":"200"}]`,
			"ip", "-detail", "-json", "route", "show", "table", "all").
		Set(`[]`, "ip", "-json", "rule").
		Set(`[{"priority":200,"src":"all","table":"200"}]`, "ip", "-6", "-json", "rule").
		Set(`[]`, "ip", "-json", "neigh", "show").
		Set(`[{"dst":"2001:db8::1","dev":"eth0","prefsrc":"2001:db8::2","table":"200"}]`,
			"ip", "-json", "-6", "route", "get", "2001:db8::1").
		Set(`[{"dst":"default","dev":"eth0","table":"200"}]`,
			"ip", "-json", "-6", "route", "get", "2001:db8::1", "fibmatch")
	collector := NewPath(fake)
	collector.Resolve = func(context.Context, string) ([]netip.Addr, error) {
		return []netip.Addr{netip.MustParseAddr("2001:db8::1")}, nil
	}

	got, warnings := collector.Collect(context.Background(), "v6.test", PathOptions{Family: "6"})
	if len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %+v", warnings)
	}
	if rule := got.Candidates[0].Rule; rule == nil || rule.Priority != 200 || rule.Family != "inet6" {
		t.Fatalf("unexpected IPv6 rule: %+v", rule)
	}
}
