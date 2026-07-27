package app

import (
	"time"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/spf13/cobra"
)

// Options holds the resolved global flags shared across commands. Command-
// specific flags (connectivity targets) are bound on their own command.
type Options struct {
	// Output
	JSON      bool
	Pretty    bool
	Summary   bool
	Barebones bool
	NoColor   bool
	Color     string // auto|always|never
	Warnings  bool
	Debug     bool

	// Filters
	Up        bool
	Interface string
	Driver    string
	State     string // up|down|unknown
	Physical  bool
	Virtual   bool
	VF        bool
	PF        bool

	// OVS
	OVS    bool
	NoOVS  bool
	NoSudo bool

	// Cross-cutting
	Watch time.Duration
	Netns string

	// Statistics (bound on the stats command).
	StatsRate       time.Duration
	StatsSort       string
	StatsTop        int
	statsPrevious   []model.InterfaceStats
	statsPreviousAt time.Time

	// Connectivity (bound on the test command)
	PingTarget       string
	DNSTarget        string
	WebTarget        string
	ThroughputTarget string
	Throughput       bool
	Count            int
	Timeout          time.Duration
}

// bindGlobal registers global persistent flags on cmd, writing into o. Old
// netcheck short flags are preserved as aliases where they map cleanly.
func (o *Options) bindGlobal(cmd *cobra.Command) {
	f := cmd.PersistentFlags()
	f.BoolVarP(&o.JSON, "json", "j", false, "emit JSON instead of tables")
	f.BoolVar(&o.Pretty, "pretty", false, "pretty-print JSON output")
	f.BoolVarP(&o.Summary, "summary", "s", false, "shorter summary tables")
	f.BoolVarP(&o.Barebones, "barebones", "b", false, "plain pipe-delimited tables")
	f.BoolVar(&o.NoColor, "no-color", false, "disable colored output")
	f.StringVar(&o.Color, "color", "auto", "color mode: auto|always|never")
	f.BoolVar(&o.Warnings, "warnings", false, "always show collection warnings")
	f.BoolVar(&o.Debug, "debug", false, "enable debug logging")

	f.BoolVarP(&o.Up, "up", "u", false, "only show interfaces that are UP")
	f.StringVar(&o.Interface, "interface", "", "filter to a single interface name")
	f.StringVar(&o.Driver, "driver", "", "filter by driver name")
	f.StringVar(&o.State, "state", "", "filter by state: up|down|unknown")
	f.BoolVar(&o.Physical, "physical", false, "only physical interfaces")
	f.BoolVar(&o.Virtual, "virtual", false, "only virtual interfaces")
	f.BoolVar(&o.VF, "vf", false, "only SR-IOV virtual functions")
	f.BoolVar(&o.PF, "pf", false, "only SR-IOV physical functions")

	f.BoolVar(&o.OVS, "ovs", false, "include Open vSwitch data")
	f.BoolVar(&o.NoOVS, "no-ovs", false, "skip Open vSwitch data")
	f.BoolVar(&o.NoSudo, "no-sudo", false, "never escalate with sudo")

	f.DurationVar(&o.Watch, "watch", 0, "refresh the view on this interval (e.g. 2s); Ctrl-C to stop")
	f.StringVar(&o.Netns, "netns", "", "run inside the named network namespace (needs root)")
}
