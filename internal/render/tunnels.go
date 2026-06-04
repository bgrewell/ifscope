package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// Tunnels renders the tunnel/overlay interface table.
func (o Options) Tunnels(w io.Writer, tunnels []model.Tunnel) {
	headers := []string{"NAME", "TYPE", "VNI", "LOCAL", "REMOTE", "PORT", "TTL", "DEV"}
	rows := make([][]string, 0, len(tunnels))
	for _, t := range tunnels {
		vni := ""
		if t.VNI != nil {
			vni = strconv.Itoa(*t.VNI)
		}
		port := ""
		if t.Port != nil {
			port = strconv.Itoa(*t.Port)
		}
		rows = append(rows, []string{
			t.Name, t.Type, vni, t.Local, t.Remote, port, t.TTL, t.Dev,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
