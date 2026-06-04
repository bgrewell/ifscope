package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// IRQ renders the NIC interrupt-affinity table.
func (o Options) IRQ(w io.Writer, irqs []model.IRQ) {
	headers := []string{"DEVICE", "IRQ", "NAME", "CPUS"}
	rows := make([][]string, 0, len(irqs))
	for _, i := range irqs {
		rows = append(rows, []string{i.Device, strconv.Itoa(i.IRQ), i.Name, i.CPUs})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
