package render

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// Options controls human table rendering. It is built from the global CLI
// flags by the app layer, keeping render free of CLI dependencies.
type Options struct {
	Summary   bool
	Barebones bool
	Color     Color
}

// write renders a table in the mode selected by opt.
func (o Options) write(w io.Writer, t Table) {
	if o.Barebones {
		t.WriteBarebones(w)
		return
	}
	t.WriteUnicode(w)
}

// Interfaces renders the interface table. Driver/firmware/speed/port columns
// are populated by later enrichment; they render empty until then.
func (o Options) Interfaces(w io.Writer, ifaces []model.Interface) {
	headers := []string{"ID", "NAME", "MAC", "STATE", "ADDRS", "DRIVER", "FIRMWARE", "BUS", "SPEED", "PORT", "ALTNAMES"}
	if o.Summary {
		headers = []string{"ID", "NAME", "MAC", "STATE", "ADDRS", "DRIVER", "BUS", "SPEED", "PORT"}
	}

	rows := make([][]string, 0, len(ifaces))
	for _, i := range ifaces {
		base := []string{
			strconv.Itoa(i.ID),
			i.Name,
			i.MAC,
			o.Color.State(i.State),
			ipv4Cell(i),
			i.Driver,
		}
		if o.Summary {
			base = append(base, busDisplay(i.Bus), i.Speed, i.Port)
		} else {
			base = append(base, i.Firmware, busDisplay(i.Bus), i.Speed, i.Port, strings.Join(i.AltNames, "\n"))
		}
		rows = append(rows, base)
	}

	o.write(w, Table{Headers: headers, Rows: rows})
}

// VLANs renders the VLAN table.
func (o Options) VLANs(w io.Writer, vlans []model.Interface) {
	headers := []string{"ID", "NAME", "PARENT", "VID", "MAC", "STATE", "ADDRS"}
	rows := make([][]string, 0, len(vlans))
	for _, v := range vlans {
		rows = append(rows, []string{
			strconv.Itoa(v.ID),
			v.Name,
			v.LinkParent,
			strconv.Itoa(v.VLANID),
			v.MAC,
			o.Color.State(v.State),
			ipv4Cell(v),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// Section prints a titled section header above a table (skipped in barebones).
func (o Options) Section(w io.Writer, title string) {
	if o.Barebones {
		fmt.Fprintf(w, "# %s\n", title)
		return
	}
	fmt.Fprintf(w, "%s\n", title)
}

// ipv4Cell renders an interface's IPv4 addresses, one per line.
func ipv4Cell(i model.Interface) string {
	var lines []string
	for _, a := range i.IPv4() {
		lines = append(lines, fmt.Sprintf("%s/%d", a.Local, a.PrefixLen))
	}
	return strings.Join(lines, "\n")
}

// busDisplay trims the common "0000:" PCI domain prefix for readability. The
// full bus address is preserved in JSON output.
func busDisplay(bus string) string {
	return strings.TrimPrefix(bus, "0000:")
}
