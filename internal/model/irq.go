package model

// IRQ is a NIC interrupt and the CPUs it is allowed to run on (its SMP
// affinity), correlated to the owning interface.
type IRQ struct {
	Device string `json:"device"`
	IRQ    int    `json:"irq"`
	Name   string `json:"name,omitempty"`
	CPUs   string `json:"cpus,omitempty"`
}
