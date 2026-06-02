package parse

import (
	"encoding/json"
	"fmt"

	"github.com/bgrewell/ifscope/internal/model"
)

// OVS column selections requested from ovs-vsctl. Parsing relies on the
// "headings" array rather than fixed positions.
const (
	OVSBridgeColumns    = "name,ports"
	OVSPortColumns      = "_uuid,name,interfaces,tag,trunks,vlan_mode,bond_mode,lacp"
	OVSInterfaceColumns = "_uuid,name,type,ofport,external_ids,options"
)

// ovsTable is the JSON envelope of `ovs-vsctl --format=json --columns=... list`.
type ovsTable struct {
	Headings []string            `json:"headings"`
	Data     [][]json.RawMessage `json:"data"`
}

// rows maps each row's heading to its raw value for keyed access.
func (t ovsTable) rows() []map[string]json.RawMessage {
	out := make([]map[string]json.RawMessage, 0, len(t.Data))
	for _, row := range t.Data {
		m := make(map[string]json.RawMessage, len(t.Headings))
		for i, h := range t.Headings {
			if i < len(row) {
				m[h] = row[i]
			}
		}
		out = append(out, m)
	}
	return out
}

func parseOVSTable(data []byte) (ovsTable, error) {
	var t ovsTable
	if err := json.Unmarshal(data, &t); err != nil {
		return ovsTable{}, err
	}
	return t, nil
}

// OVS parses the Bridge, Port, and Interface tables into a correlated model,
// resolving the uuid references between them.
func OVS(bridgeData, portData, ifaceData []byte) (model.OVS, error) {
	bridges, err := parseOVSTable(bridgeData)
	if err != nil {
		return model.OVS{}, fmt.Errorf("parse ovs bridges: %w", err)
	}
	ports, err := parseOVSTable(portData)
	if err != nil {
		return model.OVS{}, fmt.Errorf("parse ovs ports: %w", err)
	}
	ifaces, err := parseOVSTable(ifaceData)
	if err != nil {
		return model.OVS{}, fmt.Errorf("parse ovs interfaces: %w", err)
	}

	// Interfaces by uuid (carrying their row uuid is implicit: the Port table
	// references them, and Port rows we key by uuid below).
	ifByUUID := map[string]model.OVSInterface{}
	var ovsIfaces []model.OVSInterface
	for _, row := range ifaces.rows() {
		oi := model.OVSInterface{
			Name:        ovsString(row["name"]),
			Type:        ovsString(row["type"]),
			ExternalIDs: ovsMap(row["external_ids"]),
			Options:     ovsMap(row["options"]),
		}
		if n, ok := ovsInt(row["ofport"]); ok {
			oi.OFPort = &n
		}
		ovsIfaces = append(ovsIfaces, oi)
		if id := rowUUID(ifaces, row); id != "" {
			ifByUUID[id] = oi
		}
	}

	// Ports by uuid, resolving their interface names.
	type portRow struct {
		port  model.OVSPort
		ifIDs []string
	}
	portByUUID := map[string]*portRow{}
	for _, row := range ports.rows() {
		p := model.OVSPort{
			Name:     ovsString(row["name"]),
			VLANMode: ovsString(row["vlan_mode"]),
			BondMode: ovsString(row["bond_mode"]),
			LACP:     ovsString(row["lacp"]),
			Trunks:   ovsInts(row["trunks"]),
		}
		if tag, ok := ovsInt(row["tag"]); ok {
			p.Tag = &tag
		}
		ifIDs := ovsUUIDs(row["interfaces"])
		for _, id := range ifIDs {
			if oi, ok := ifByUUID[id]; ok {
				p.Interfaces = append(p.Interfaces, oi.Name)
			}
		}
		pr := &portRow{port: p, ifIDs: ifIDs}
		if id := rowUUID(ports, row); id != "" {
			portByUUID[id] = pr
		}
	}

	// Bridges resolve their port names and back-fill each port's bridge.
	var ovsBridges []model.OVSBridge
	var ovsPorts []model.OVSPort
	seenPort := map[string]bool{}
	for _, row := range bridges.rows() {
		b := model.OVSBridge{Name: ovsString(row["name"])}
		for _, pid := range ovsUUIDs(row["ports"]) {
			pr, ok := portByUUID[pid]
			if !ok {
				continue
			}
			b.Ports = append(b.Ports, pr.port.Name)
			if !seenPort[pid] {
				pr.port.Bridge = b.Name
				ovsPorts = append(ovsPorts, pr.port)
				seenPort[pid] = true
			}
		}
		ovsBridges = append(ovsBridges, b)
	}

	return model.OVS{Bridges: ovsBridges, Ports: ovsPorts, Interfaces: ovsIfaces}, nil
}

// rowUUID returns the row's own _uuid when present. ovs-vsctl includes _uuid in
// --columns output only when requested; when absent this returns "" and uuid
// resolution falls back to positional set membership, which our column choices
// avoid needing. To support resolution we request _uuid implicitly via the
// "uuid" heading if present.
func rowUUID(t ovsTable, row map[string]json.RawMessage) string {
	if v, ok := row["_uuid"]; ok {
		ids := ovsUUIDs(v)
		if len(ids) == 1 {
			return ids[0]
		}
	}
	return ""
}

// ovsString decodes a scalar string OVS value, returning "" for empty sets.
func ovsString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}
	return ""
}

// ovsInt decodes a scalar integer OVS value. An empty set yields (0, false).
func ovsInt(raw json.RawMessage) (int, bool) {
	if len(raw) == 0 {
		return 0, false
	}
	var n int
	if json.Unmarshal(raw, &n) == nil {
		return n, true
	}
	return 0, false
}

// ovsUUIDs decodes a uuid, a set of uuids, or an empty set into uuid strings.
// OVS omits the "set" wrapper for single-element sets, so both shapes are
// handled.
func ovsUUIDs(raw json.RawMessage) []string {
	var arr []json.RawMessage
	if len(raw) == 0 || json.Unmarshal(raw, &arr) != nil || len(arr) != 2 {
		return nil
	}
	var tag string
	if json.Unmarshal(arr[0], &tag) != nil {
		return nil
	}
	switch tag {
	case "uuid":
		var id string
		if json.Unmarshal(arr[1], &id) == nil {
			return []string{id}
		}
	case "set":
		var elems []json.RawMessage
		if json.Unmarshal(arr[1], &elems) == nil {
			var out []string
			for _, e := range elems {
				out = append(out, ovsUUIDs(e)...)
			}
			return out
		}
	}
	return nil
}

// ovsInts decodes an integer, a set of integers, or an empty set.
func ovsInts(raw json.RawMessage) []int {
	if n, ok := ovsInt(raw); ok {
		return []int{n}
	}
	var arr []json.RawMessage
	if len(raw) == 0 || json.Unmarshal(raw, &arr) != nil || len(arr) != 2 {
		return nil
	}
	var tag string
	if json.Unmarshal(arr[0], &tag) != nil || tag != "set" {
		return nil
	}
	var elems []int
	if json.Unmarshal(arr[1], &elems) == nil {
		return elems
	}
	return nil
}

// ovsMap decodes an OVS map value into a string→string map.
func ovsMap(raw json.RawMessage) map[string]string {
	var arr []json.RawMessage
	if len(raw) == 0 || json.Unmarshal(raw, &arr) != nil || len(arr) != 2 {
		return nil
	}
	var tag string
	if json.Unmarshal(arr[0], &tag) != nil || tag != "map" {
		return nil
	}
	var pairs [][]string
	if json.Unmarshal(arr[1], &pairs) != nil {
		return nil
	}
	m := make(map[string]string, len(pairs))
	for _, p := range pairs {
		if len(p) == 2 {
			m[p[0]] = p[1]
		}
	}
	if len(m) == 0 {
		return nil
	}
	return m
}
