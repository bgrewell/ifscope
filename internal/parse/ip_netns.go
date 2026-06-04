package parse

import (
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// NetnsList parses `ip netns list` text. Each line is a namespace name,
// optionally followed by "(id: N)", e.g.:
//
//	myns (id: 0)
//	other
func NetnsList(data []byte) []model.Netns {
	var out []model.Netns
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		ns := model.Netns{Name: fields[0]}
		if i := strings.Index(line, "id:"); i >= 0 {
			rest := strings.TrimRight(strings.TrimSpace(line[i+len("id:"):]), ")")
			if id, err := strconv.Atoi(strings.TrimSpace(rest)); err == nil {
				ns.ID = &id
			}
		}
		out = append(out, ns)
	}
	return out
}
