package parse

import (
	"regexp"
	"strings"
)

// PCIInfo is the device description and vendor/device IDs extracted from one
// `lspci -Dnn -s <bus>` line.
type PCIInfo struct {
	Description string
	VendorID    string
	DeviceID    string
}

// idRe matches the "[vendor:device]" hex pair lspci appends in -nn mode.
var idRe = regexp.MustCompile(`\[([0-9a-fA-F]{4}):([0-9a-fA-F]{4})\]`)

// LspciDevice parses a single `lspci -Dnn -s <bus>` line, e.g.:
//
//	0000:17:00.0 Ethernet controller [0200]: Intel Corporation E810 [8086:1592] (rev 02)
//
// It returns the vendor/device name as Description plus the vendor and device
// IDs. An empty or unrecognized line yields a zero PCIInfo.
func LspciDevice(data []byte) PCIInfo {
	line := strings.TrimSpace(firstLine(string(data)))
	if line == "" {
		return PCIInfo{}
	}

	var info PCIInfo
	if m := idRe.FindStringSubmatch(line); m != nil {
		info.VendorID, info.DeviceID = m[1], m[2]
	}

	// The vendor/device name follows the first ": " (after "<bus> <class> [id]:").
	_, rest, ok := strings.Cut(line, ": ")
	if !ok {
		return info
	}
	// Strip the trailing "[vendor:device]" and any "(rev ..)" suffix.
	rest = idRe.ReplaceAllString(rest, "")
	if i := strings.LastIndex(rest, "(rev "); i >= 0 {
		rest = rest[:i]
	}
	info.Description = strings.TrimSpace(rest)
	return info
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
