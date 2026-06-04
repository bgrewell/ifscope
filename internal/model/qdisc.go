package model

// Qdisc is a queueing discipline attached to an interface. ifscope reports the
// root qdisc per device (the effective discipline, e.g. mq, fq_codel, htb);
// child classes/filters are out of scope for this view.
type Qdisc struct {
	Dev    string `json:"dev"`
	Kind   string `json:"kind"`
	Handle string `json:"handle,omitempty"`
	Parent string `json:"parent,omitempty"`
	Root   bool   `json:"root"`
}
