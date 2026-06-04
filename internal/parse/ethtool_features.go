package parse

import "strings"

// EthtoolFeatures parses `ethtool -k <iface>` output into top-level offload
// features. Indented sub-features and the "Features for X:" header are skipped.
// Values are kept verbatim (e.g. "on", "off", "off [fixed]").
func EthtoolFeatures(data []byte) map[string]string {
	out := map[string]string{}
	for _, raw := range strings.Split(string(data), "\n") {
		// Top-level features are not indented; sub-features start with a tab/space.
		if raw == "" || raw[0] == '\t' || raw[0] == ' ' {
			continue
		}
		key, val, ok := strings.Cut(raw, ": ")
		if !ok {
			continue // header line "Features for eth0:" has no ": value"
		}
		out[strings.TrimSpace(key)] = strings.TrimSpace(val)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
