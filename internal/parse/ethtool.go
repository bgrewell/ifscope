package parse

import (
	"strconv"
	"strings"
)

// DriverInfo is the subset of `ethtool -i <iface>` ifscope consumes.
type DriverInfo struct {
	Driver   string
	Version  string
	Firmware string
	Bus      string
}

// LinkSettings is the subset of `ethtool <iface>` ifscope consumes, with Speed
// normalized to a human-friendly form.
type LinkSettings struct {
	Speed  string
	Port   string
	Duplex string
}

// EthtoolDriverInfo parses `ethtool -i` output. Parsing is tolerant: only lines
// of the form "key: value" are considered and unknown keys are ignored.
func EthtoolDriverInfo(data []byte) DriverInfo {
	kv := parseColonKV(data)
	return DriverInfo{
		Driver:   kv["driver"],
		Version:  kv["version"],
		Firmware: normalizeFirmware(kv["firmware-version"]),
		Bus:      kv["bus-info"],
	}
}

// EthtoolSettings parses `ethtool <iface>` output. Continuation lines (e.g. the
// multi-line supported/advertised link-mode lists) and noise lines such as
// "netlink error: ..." are ignored; only the fields of interest are extracted.
func EthtoolSettings(data []byte) LinkSettings {
	var s LinkSettings
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		key, val, ok := strings.Cut(line, ": ")
		if !ok {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Speed":
			s.Speed = normalizeSpeed(strings.TrimSpace(val))
		case "Port":
			s.Port = strings.TrimSpace(val)
		case "Duplex":
			s.Duplex = strings.TrimSpace(val)
		}
	}
	return s
}

// parseColonKV builds a key→value map from "key: value" lines, lowercasing
// keys and keeping the first occurrence of each.
func parseColonKV(data []byte) map[string]string {
	out := map[string]string{}
	for _, raw := range strings.Split(string(data), "\n") {
		key, val, ok := strings.Cut(raw, ":")
		if !ok {
			continue
		}
		k := strings.ToLower(strings.TrimSpace(key))
		if _, seen := out[k]; seen {
			continue
		}
		out[k] = strings.TrimSpace(val)
	}
	return out
}

// normalizeFirmware drops a placeholder firmware value.
func normalizeFirmware(s string) string {
	if s == "N/A" || s == "" {
		return ""
	}
	return s
}

// normalizeSpeed converts ethtool's "<n>Mb/s" into Gb/s when it divides evenly,
// otherwise keeps Mb/s. Unknown speeds normalize to the empty string.
func normalizeSpeed(s string) string {
	if s == "" || strings.HasPrefix(strings.ToLower(s), "unknown") {
		return ""
	}
	num := strings.TrimSuffix(s, "Mb/s")
	if num == s {
		return s // unexpected unit; pass through
	}
	mbps, err := strconv.Atoi(strings.TrimSpace(num))
	if err != nil {
		return s
	}
	if mbps >= 1000 && mbps%1000 == 0 {
		return strconv.Itoa(mbps/1000) + " Gb/s"
	}
	return strconv.Itoa(mbps) + " Mb/s"
}
