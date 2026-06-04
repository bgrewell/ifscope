package render

import (
	"io"

	"github.com/bgrewell/ifscope/internal/model"
)

// Sockets renders the listening-socket table.
func (o Options) Sockets(w io.Writer, sockets []model.Socket) {
	headers := []string{"PROTO", "STATE", "LOCAL", "PORT", "PROCESS"}
	rows := make([][]string, 0, len(sockets))
	for _, s := range sockets {
		rows = append(rows, []string{
			s.Proto,
			s.State,
			s.LocalAddr,
			s.LocalPort,
			s.Process,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
