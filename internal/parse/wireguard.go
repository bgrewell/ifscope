package parse

import (
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// WireGuard parses `wg show all dump`. Lines are tab-separated; an interface
// line has 5 fields (iface, private-key, public-key, listen-port, fwmark) and a
// peer line has 9 (iface, public-key, preshared-key, endpoint, allowed-ips,
// latest-handshake, rx, tx, keepalive). Private and pre-shared keys are
// deliberately discarded.
func WireGuard(data []byte) []model.WireGuardDevice {
	var out []model.WireGuardDevice
	idx := map[string]int{} // iface -> index in out

	for _, line := range strings.Split(string(data), "\n") {
		if line == "" {
			continue
		}
		f := strings.Split(line, "\t")
		switch len(f) {
		case 5: // interface line (private key in f[1] is ignored)
			port, _ := strconv.Atoi(f[3])
			idx[f[0]] = len(out)
			out = append(out, model.WireGuardDevice{
				Interface:  f[0],
				PublicKey:  f[2],
				ListenPort: port,
			})
		case 9: // peer line (pre-shared key in f[2] is ignored)
			i, ok := idx[f[0]]
			if !ok {
				continue
			}
			out[i].Peers = append(out[i].Peers, parseWGPeer(f))
		}
	}
	return out
}

func parseWGPeer(f []string) model.WGPeer {
	p := model.WGPeer{PublicKey: f[1]}
	if f[3] != "(none)" {
		p.Endpoint = f[3]
	}
	if f[4] != "(none)" {
		p.AllowedIPs = strings.Split(f[4], ",")
	}
	if hs, err := strconv.ParseInt(f[5], 10, 64); err == nil {
		p.LatestHandshake = hs
	}
	if rx, err := strconv.ParseUint(f[6], 10, 64); err == nil {
		p.RxBytes = rx
	}
	if tx, err := strconv.ParseUint(f[7], 10, 64); err == nil {
		p.TxBytes = tx
	}
	if f[8] != "off" {
		p.Keepalive = f[8]
	}
	return p
}
