package collect

import (
	"context"
	"regexp"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

const ethtoolCmd = "ethtool"

// pciBusRe matches a fully-qualified PCI address (domain:bus:device.function).
var pciBusRe = regexp.MustCompile(`^[0-9a-fA-F]{4}:[0-9a-fA-F]{2}:[0-9a-fA-F]{2}\.[0-9a-fA-F]+$`)

// isPCIBus reports whether bus is a PCI address ifscope can map to sysfs/lspci.
func isPCIBus(bus string) bool { return pciBusRe.MatchString(bus) }

// Ethtool enriches interfaces with driver, firmware, speed, and port details.
type Ethtool struct {
	Runner run.Runner
}

// NewEthtool returns an Ethtool collector using r.
func NewEthtool(r run.Runner) *Ethtool { return &Ethtool{Runner: r} }

// Enrich fills driver/firmware (from `ethtool -i`) and speed/port (from
// `ethtool <iface>`) in place. Existing values are preserved. Per-interface
// failures (common on virtual devices that lack link settings) are ignored; a
// missing ethtool binary yields a single non-fatal warning.
func (c *Ethtool) Enrich(ctx context.Context, ifaces []model.Interface) []model.Warning {
	for i := range ifaces {
		name := ifaces[i].Name
		if name == "" {
			continue
		}

		out, _, err := c.Runner.Run(ctx, ethtoolCmd, "-i", name)
		if err != nil {
			if run.IsNotFound(err) {
				return []model.Warning{{
					Source:  "ethtool",
					Message: "ethtool not found; driver, firmware, speed, and port details unavailable",
				}}
			}
			// Some interfaces (e.g. certain virtual devices) reject -i; skip.
		} else {
			di := parse.EthtoolDriverInfo(out)
			if ifaces[i].Driver == "" {
				ifaces[i].Driver = di.Driver
			}
			if ifaces[i].Firmware == "" {
				ifaces[i].Firmware = di.Firmware
			}
			if ifaces[i].Bus == "" && isPCIBus(di.Bus) {
				ifaces[i].Bus = di.Bus
			}
		}

		// Link settings are meaningful only for ports that have them; ignore
		// the "Operation not supported" failures on virtual interfaces.
		sout, _, serr := c.Runner.Run(ctx, ethtoolCmd, name)
		if serr == nil {
			ls := parse.EthtoolSettings(sout)
			if ifaces[i].Speed == "" {
				ifaces[i].Speed = ls.Speed
			}
			if ifaces[i].Port == "" {
				ifaces[i].Port = ls.Port
			}
		}
	}
	return nil
}
