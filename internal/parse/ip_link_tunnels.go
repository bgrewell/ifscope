package parse

import (
	"bufio"
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// tunnelKinds are the link kinds the tunnels view understands.
var tunnelKinds = map[string]bool{
	"vxlan": true, "geneve": true, "gre": true, "gretap": true,
	"ip6gre": true, "erspan": true, "sit": true, "ipip": true,
	"ip6tnl": true, "vti": true, "vti6": true,
}

// tunnelHeaderRe matches an interface header line, e.g. "7: gnv200: <...>".
var tunnelHeaderRe = regexp.MustCompile(`^\d+:\s+([^:@\s]+)`)

// Tunnels parses `ip -d link show` text into tunnel interfaces. Text is used
// rather than `-json` because iproute2 emits malformed JSON for vxlan (a stray
// "fan-map" token), whereas the text detail line is reliable across versions.
func Tunnels(data []byte) []model.Tunnel {
	var out []model.Tunnel
	var current string

	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if m := tunnelHeaderRe.FindStringSubmatch(line); m != nil {
			current = m[1]
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 || !tunnelKinds[fields[0]] {
			continue
		}
		out = append(out, parseTunnelDetail(current, fields))
	}
	return out
}

// parseTunnelDetail walks the key/value tokens of a tunnel detail line. Flag
// tokens without a value (e.g. "fan-map", "udpcsum") are skipped.
func parseTunnelDetail(name string, fields []string) model.Tunnel {
	t := model.Tunnel{Name: name, Type: fields[0]}
	for i := 1; i < len(fields)-1; i++ {
		switch fields[i] {
		case "id", "vni":
			if v, err := strconv.Atoi(fields[i+1]); err == nil {
				t.VNI = &v
			}
		case "remote":
			t.Remote = fields[i+1]
		case "local":
			t.Local = fields[i+1]
		case "dstport":
			if v, err := strconv.Atoi(fields[i+1]); err == nil {
				t.Port = &v
			}
		case "ttl":
			t.TTL = fields[i+1]
		case "dev":
			t.Dev = fields[i+1]
		}
	}
	return t
}
