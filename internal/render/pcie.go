package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// PCIe renders the PCIe device table.
func (o Options) PCIe(w io.Writer, devices []model.PCIDevice) {
	headers := []string{"BUS", "INTERFACE", "DRIVER", "BIND", "DEVICE", "VENDOR", "NUMA", "LINK"}
	rows := make([][]string, 0, len(devices))
	for _, d := range devices {
		rows = append(rows, []string{
			busDisplay(d.Bus),
			d.Interface,
			d.Driver,
			d.Bind,
			d.Description,
			vendorCell(d),
			numaCell(d.NUMANode),
			linkCell(d),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// vendorCell shows the vendor:device id pair when known.
func vendorCell(d model.PCIDevice) string {
	if d.VendorID == "" && d.DeviceID == "" {
		return ""
	}
	return d.VendorID + ":" + d.DeviceID
}

// numaCell renders an optional NUMA node, blank when unknown.
func numaCell(node *int) string {
	if node == nil {
		return ""
	}
	return strconv.Itoa(*node)
}

// linkCell combines PCIe link speed and width when available.
func linkCell(d model.PCIDevice) string {
	switch {
	case d.LinkSpeed != "" && d.LinkWidth != "":
		return d.LinkSpeed + " x" + d.LinkWidth
	case d.LinkSpeed != "":
		return d.LinkSpeed
	default:
		return ""
	}
}
