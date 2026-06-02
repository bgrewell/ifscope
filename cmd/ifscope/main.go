// Command ifscope is a Linux CLI for inspecting host network state — interfaces,
// VLANs, routes, DNS, PCIe NIC mapping, OVS topology, SR-IOV/VF state, and basic
// connectivity tests — in human tables or stable JSON.
package main

import (
	"os"

	"github.com/bgrewell/ifscope/internal/app"
)

func main() {
	os.Exit(app.Execute())
}
