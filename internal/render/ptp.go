package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// PTP renders the hardware-timestamping / PTP capability table.
func (o Options) PTP(w io.Writer, ptps []model.PTP) {
	headers := []string{"NAME", "PHC", "HW-TX", "HW-RX", "SW-TX", "SW-RX"}
	rows := make([][]string, 0, len(ptps))
	for _, p := range ptps {
		phc := ""
		if p.PHCIndex != nil {
			phc = strconv.Itoa(*p.PHCIndex)
		}
		rows = append(rows, []string{
			p.Name, phc, yesNo(p.HWTx), yesNo(p.HWRx), yesNo(p.SWTx), yesNo(p.SWRx),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
