package render

import (
	"io"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// Routes renders the routing table.
func (o Options) Routes(w io.Writer, routes []model.Route) {
	headers := []string{"DST", "GATEWAY", "DEV", "PROTO", "METRIC", "TABLE", "SCOPE", "SRC"}
	rows := make([][]string, 0, len(routes))
	for _, r := range routes {
		metric := ""
		if r.Metric != 0 {
			metric = strconv.Itoa(r.Metric)
		}
		rows = append(rows, []string{
			r.Dst, r.Gateway, r.Dev, r.Protocol, metric, r.Table, r.Scope, r.Src,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// DNS renders the per-link DNS table.
func (o Options) DNS(w io.Writer, dns []model.DNS) {
	headers := []string{"LINK", "CURRENT", "SERVERS", "DOMAINS", "DEFROUTE", "LLMNR", "MDNS", "DNSSEC"}
	rows := make([][]string, 0, len(dns))
	for _, d := range dns {
		rows = append(rows, []string{
			d.Link,
			d.CurrentServer,
			strings.Join(d.Servers, "\n"),
			strings.Join(d.Domains, "\n"),
			boolCell(d.DefaultRoute),
			d.LLMNR,
			d.MDNS,
			d.DNSSEC,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// boolCell renders an optional bool as yes/no, blank when unset.
func boolCell(b *bool) string {
	if b == nil {
		return ""
	}
	if *b {
		return "yes"
	}
	return "no"
}
