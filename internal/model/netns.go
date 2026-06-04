package model

// Netns is a named network namespace. Interfaces is the count of links inside
// it, or nil when that could not be determined (e.g. without privileges).
type Netns struct {
	Name       string `json:"name"`
	ID         *int   `json:"id,omitempty"`
	Interfaces *int   `json:"interfaces,omitempty"`
}
