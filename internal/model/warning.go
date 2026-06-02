// Package model defines the pure data types that make up an ifscope report.
//
// Types here carry no I/O behavior; collectors populate them and renderers
// consume them. JSON tags define the stable machine-readable schema.
package model

// Warning records a non-fatal issue encountered during collection.
//
// Collectors return warnings instead of printing or aborting, so that missing
// optional dependencies or permission failures degrade gracefully. A warning
// with Fatal set means the requested operation could not continue.
type Warning struct {
	Source  string `json:"source"`
	Message string `json:"message"`
	Fatal   bool   `json:"fatal,omitempty"`
}
