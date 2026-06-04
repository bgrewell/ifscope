package parse

import (
	"regexp"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// ssProcessRe extracts ("name",pid=NNN) tuples from ss's process column.
var ssProcessRe = regexp.MustCompile(`"([^"]+)",pid=(\d+)`)

// SS parses headerless `ss -tulpn` output into listening sockets. Columns are
// whitespace-separated: netid state recv-q send-q local:port peer:port [users].
func SS(data []byte) []model.Socket {
	var out []model.Socket
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}
		addr, port := splitLastColon(fields[4])
		s := model.Socket{
			Proto:     fields[0],
			State:     fields[1],
			LocalAddr: addr,
			LocalPort: port,
		}
		if len(fields) >= 7 {
			s.Process = ssProcesses(strings.Join(fields[6:], " "))
		}
		out = append(out, s)
	}
	return out
}

// splitLastColon splits "host:port" at the final colon, handling IPv6 forms
// like "[::]:22" and "127.0.0.53%lo:53".
func splitLastColon(s string) (host, port string) {
	i := strings.LastIndex(s, ":")
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

// ssProcesses renders the process column as "name(pid)" entries.
func ssProcesses(s string) string {
	matches := ssProcessRe.FindAllStringSubmatch(s, -1)
	if matches == nil {
		return ""
	}
	parts := make([]string, 0, len(matches))
	for _, m := range matches {
		parts = append(parts, m[1]+"("+m[2]+")")
	}
	return strings.Join(parts, ",")
}
