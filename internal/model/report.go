package model

// Host carries basic identity of the inspected host.
type Host struct {
	Hostname string `json:"hostname,omitempty"`
}

// Report is the top-level correlated result produced by a run. Renderers emit
// it as a human table view or as JSON. The Version field versions the JSON
// schema so consumers can adapt to future changes.
type Report struct {
	Version    string       `json:"version"`
	Host       Host         `json:"host"`
	Interfaces []Interface  `json:"interfaces,omitempty"`
	VLANs      []Interface  `json:"vlans,omitempty"`
	Bonds      []Bond       `json:"bonds,omitempty"`
	Routes     []Route      `json:"routes,omitempty"`
	Rules      []Rule       `json:"rules,omitempty"`
	DNS        []DNS        `json:"dns,omitempty"`
	PCIe       []PCIDevice  `json:"pcie,omitempty"`
	OVS        *OVS         `json:"ovs,omitempty"`
	Tests      []TestResult `json:"tests,omitempty"`
	Warnings   []Warning    `json:"warnings,omitempty"`
}

// AddWarning appends a warning to the report.
func (r *Report) AddWarning(w Warning) {
	r.Warnings = append(r.Warnings, w)
}
