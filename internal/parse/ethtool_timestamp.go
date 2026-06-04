package parse

import (
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// EthtoolTimestamp parses `ethtool -T` into PTP/timestamping capabilities. The
// returned PTP has no Name (the collector sets it).
func EthtoolTimestamp(data []byte) model.PTP {
	var p model.PTP
	inCaps := false
	for _, raw := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(raw)
		indented := raw != trimmed

		if trimmed == "Capabilities:" {
			inCaps = true
			continue
		}
		if inCaps && indented {
			switch trimmed {
			case "hardware-transmit":
				p.HWTx = true
			case "hardware-receive":
				p.HWRx = true
			case "software-transmit":
				p.SWTx = true
			case "software-receive":
				p.SWRx = true
			}
			continue
		}
		inCaps = false

		// PHC index: newer ethtool uses "Hardware timestamp provider index",
		// older uses "PTP Hardware Clock".
		if v, ok := afterColon(trimmed, "Hardware timestamp provider index"); ok {
			if n, err := strconv.Atoi(v); err == nil {
				p.PHCIndex = &n
			}
		} else if v, ok := afterColon(trimmed, "PTP Hardware Clock"); ok {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				p.PHCIndex = &n
			}
		}
	}
	return p
}

// afterColon returns the trimmed value when line is "<prefix>: <value>".
func afterColon(line, prefix string) (string, bool) {
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(line, prefix)
	rest = strings.TrimSpace(rest)
	if !strings.HasPrefix(rest, ":") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(rest, ":")), true
}
