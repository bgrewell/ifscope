package render

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/bgrewell/ifscope/internal/model"
)

// WireGuard renders the WireGuard table: one row per peer, or a single row for
// a peerless interface.
func (o Options) WireGuard(w io.Writer, devices []model.WireGuardDevice) {
	headers := []string{"IFACE", "LISTEN", "PEER", "ENDPOINT", "ALLOWED IPS", "HANDSHAKE", "RX", "TX"}
	var rows [][]string
	for _, d := range devices {
		listen := ""
		if d.ListenPort > 0 {
			listen = strconv.Itoa(d.ListenPort)
		}
		if len(d.Peers) == 0 {
			rows = append(rows, []string{d.Interface, listen, "", "", "", "", "", ""})
			continue
		}
		for _, p := range d.Peers {
			rows = append(rows, []string{
				d.Interface,
				listen,
				p.PublicKey,
				p.Endpoint,
				strings.Join(p.AllowedIPs, "\n"),
				handshakeCell(p.LatestHandshake),
				humanBytes(p.RxBytes),
				humanBytes(p.TxBytes),
			})
		}
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// handshakeCell renders a latest-handshake epoch as a relative age, or "never".
func handshakeCell(epoch int64) string {
	if epoch <= 0 {
		return "never"
	}
	d := time.Since(time.Unix(epoch, 0)).Round(time.Second)
	if d < 0 {
		d = 0
	}
	return d.String() + " ago"
}
