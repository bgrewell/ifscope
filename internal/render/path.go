package render

import (
	"io"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// Path renders one row per resolved destination address.
func (o Options) Path(w io.Writer, path model.Path) {
	headers := []string{"ADDRESS", "FAMILY", "TABLE/RULE", "SOURCE", "VIA", "DEV", "NEIGHBOR", "TOPOLOGY", "MTU", "ERROR"}
	rows := make([][]string, 0, len(path.Candidates))
	for _, c := range path.Candidates {
		var tableRule, source, via, dev, neighbor, topology, mtu string
		if c.Route != nil {
			tableRule = c.Route.Table
			source, via, dev = c.Route.Src, c.Route.Gateway, c.Route.Dev
		}
		if c.Rule != nil {
			tableRule += "/" + strconv.Itoa(c.Rule.Priority)
		}
		if c.Neighbor != nil {
			neighbor = c.Neighbor.State
			if c.Neighbor.LLAddr != "" {
				neighbor += " " + c.Neighbor.LLAddr
			}
		}
		names := make([]string, 0, len(c.Topology))
		for _, hop := range c.Topology {
			names = append(names, hop.Name+" ("+string(hop.Type)+")")
		}
		topology = strings.Join(names, " → ")
		if c.MTU != 0 {
			mtu = strconv.Itoa(c.MTU)
		}
		rows = append(rows, []string{c.Address, c.Family, tableRule, source, via, dev, neighbor, topology, mtu, c.Error})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
