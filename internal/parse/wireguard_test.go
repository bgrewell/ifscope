package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestWireGuard(t *testing.T) {
	// Fixture captured from `wg show all dump` in a live container (private key
	// redacted; the parser discards it regardless).
	devs := WireGuard(testutil.Fixture(t, "wireguard/dump.txt"))
	if len(devs) != 1 {
		t.Fatalf("devices = %d, want 1", len(devs))
	}
	d := devs[0]
	if d.Interface != "wg0" || d.ListenPort != 51820 {
		t.Errorf("device = %+v", d)
	}
	if d.PublicKey != "SLNuQHeIqBU6aTpjZIFnG6CWa1oZax0FYMjWeIh1xiY=" {
		t.Errorf("public key = %q", d.PublicKey)
	}
	if len(d.Peers) != 1 {
		t.Fatalf("peers = %d, want 1", len(d.Peers))
	}
	p := d.Peers[0]
	if p.Endpoint != "192.0.2.5:51820" {
		t.Errorf("endpoint = %q", p.Endpoint)
	}
	if len(p.AllowedIPs) != 2 || p.AllowedIPs[1] != "fd00::/64" {
		t.Errorf("allowed-ips = %v", p.AllowedIPs)
	}
	if p.TxBytes != 148 || p.Keepalive != "25" || p.LatestHandshake != 0 {
		t.Errorf("peer stats = %+v", p)
	}
}

func TestWireGuardEmpty(t *testing.T) {
	if got := WireGuard(nil); len(got) != 0 {
		t.Errorf("empty = %v", got)
	}
}
