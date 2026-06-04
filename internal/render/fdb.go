package render

import (
	"io"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// FDB renders the bridge forwarding-database table.
func (o Options) FDB(w io.Writer, entries []model.FDBEntry) {
	headers := []string{"MAC", "DEV", "VLAN", "MASTER", "STATE", "FLAGS"}
	rows := make([][]string, 0, len(entries))
	for _, e := range entries {
		vlan := ""
		if e.VLAN != nil {
			vlan = strconv.Itoa(*e.VLAN)
		}
		rows = append(rows, []string{
			e.MAC,
			e.Dev,
			vlan,
			e.Master,
			e.State,
			strings.Join(e.Flags, ","),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
