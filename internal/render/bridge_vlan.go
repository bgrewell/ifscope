package render

import (
	"io"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// BridgeVLANs renders the bridge VLAN-filtering table.
func (o Options) BridgeVLANs(w io.Writer, vlans []model.BridgeVLAN) {
	headers := []string{"PORT", "VLAN", "FLAGS"}
	rows := make([][]string, 0, len(vlans))
	for _, v := range vlans {
		vlan := strconv.Itoa(v.VLAN)
		if v.VLANEnd > 0 {
			vlan += "-" + strconv.Itoa(v.VLANEnd)
		}
		rows = append(rows, []string{v.Port, vlan, strings.Join(v.Flags, ",")})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
