package parse

import (
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// knownDNSKeys are the resolvectl field labels ifscope recognizes. Detecting
// field lines by known key (rather than "contains a colon") is necessary
// because wrapped IPv6 server entries also contain colons.
var knownDNSKeys = map[string]bool{
	"Current Scopes":     true,
	"Protocols":          true,
	"Current DNS Server": true,
	"DNS Servers":        true,
	"DNS Domain":         true,
	"Default Route":      true,
	"resolv.conf mode":   true,
}

// Resolvectl parses `resolvectl status` output into per-link DNS state. The
// global scope is returned with Link == "global".
func Resolvectl(data []byte) []model.DNS {
	var out []model.DNS
	var cur *model.DNS
	var pendingList *[]string // continuation target for wrapped server/domain lists

	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimRight(raw, "\r")
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		switch {
		case trimmed == "Global":
			out = append(out, model.DNS{Link: "global"})
			cur = &out[len(out)-1]
			pendingList = nil
			continue
		case strings.HasPrefix(trimmed, "Link "):
			out = append(out, model.DNS{Link: linkName(trimmed)})
			cur = &out[len(out)-1]
			pendingList = nil
			continue
		}

		if cur == nil {
			continue
		}

		key, val, ok := splitDNSField(trimmed)
		if ok && knownDNSKeys[key] {
			pendingList = applyDNSField(cur, key, val)
			continue
		}

		// A non-field line continues the most recent wrapped list value.
		if pendingList != nil {
			*pendingList = append(*pendingList, strings.Fields(trimmed)...)
		}
	}
	return out
}

// applyDNSField sets the field on cur and returns a pointer to the list slice
// when the field is a wrappable list (DNS Servers / DNS Domain), else nil.
func applyDNSField(cur *model.DNS, key, val string) *[]string {
	switch key {
	case "Current DNS Server":
		cur.CurrentServer = val
	case "DNS Servers":
		cur.Servers = append(cur.Servers, strings.Fields(val)...)
		return &cur.Servers
	case "DNS Domain":
		cur.Domains = append(cur.Domains, strings.Fields(val)...)
		return &cur.Domains
	case "Default Route":
		b := val == "yes"
		cur.DefaultRoute = &b
	case "Protocols":
		applyProtocols(cur, val)
	}
	return nil
}

// applyProtocols decodes the +/- protocol flags and DNSSEC mode.
func applyProtocols(cur *model.DNS, val string) {
	for _, tok := range strings.Fields(val) {
		switch {
		case tok == "+LLMNR":
			cur.LLMNR = "yes"
		case tok == "-LLMNR":
			cur.LLMNR = "no"
		case tok == "+mDNS":
			cur.MDNS = "yes"
		case tok == "-mDNS":
			cur.MDNS = "no"
		case strings.HasPrefix(tok, "DNSSEC="):
			cur.DNSSEC = strings.TrimPrefix(tok, "DNSSEC=")
		}
	}
}

// splitDNSField splits "Key: value" on the first ": ", returning the trimmed
// key and value.
func splitDNSField(line string) (key, val string, ok bool) {
	k, v, found := strings.Cut(line, ":")
	if !found {
		return "", "", false
	}
	return strings.TrimSpace(k), strings.TrimSpace(v), true
}

// linkName extracts "eth0" from "Link 2 (eth0)", falling back to the raw text.
func linkName(s string) string {
	open := strings.IndexByte(s, '(')
	close := strings.IndexByte(s, ')')
	if open >= 0 && close > open {
		return s[open+1 : close]
	}
	return strings.TrimPrefix(s, "Link ")
}
