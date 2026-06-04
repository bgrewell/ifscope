package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// Bridges renders the Linux bridge table: one row per bridge, members one per
// line with their link state.
func (o Options) Bridges(w io.Writer, bridges []model.Bridge) {
	headers := []string{"BRIDGE", "STATE", "STP", "MEMBERS"}
	rows := make([][]string, 0, len(bridges))
	for _, b := range bridges {
		members := make([]string, 0, len(b.Members))
		for _, m := range b.Members {
			members = append(members, fmt.Sprintf("%s (%s)", m.Name, strings.ToLower(m.State)))
		}
		rows = append(rows, []string{
			b.Name,
			o.Color.State(b.State),
			yesNo(b.STP),
			strings.Join(members, "\n"),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// yesNo renders a bool as yes/no.
func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
