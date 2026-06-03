package render

import (
	"fmt"
	"io"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// Bonds renders the bonding table: one row per bond, members listed one per
// line with their link state.
func (o Options) Bonds(w io.Writer, bonds []model.Bond) {
	headers := []string{"BOND", "MODE", "STATE", "ACTIVE", "MEMBERS"}
	rows := make([][]string, 0, len(bonds))
	for _, b := range bonds {
		members := make([]string, 0, len(b.Members))
		for _, m := range b.Members {
			members = append(members, fmt.Sprintf("%s (%s)", m.Name, strings.ToLower(m.State)))
		}
		active := b.ActiveSlave
		if active == "" {
			active = "-"
		}
		rows = append(rows, []string{
			b.Name,
			b.Mode,
			o.Color.State(b.State),
			active,
			strings.Join(members, "\n"),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
