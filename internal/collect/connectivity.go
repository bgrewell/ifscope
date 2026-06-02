package collect

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/parse"
	"github.com/bgrewell/ifscope/internal/run"
)

// Default connectivity test targets. They are public, stable endpoints; all are
// overridable via flags.
const (
	DefaultPingTarget       = "1.1.1.1"
	DefaultDNSTarget        = "www.cloudflare.com"
	DefaultWebTarget        = "https://www.cloudflare.com/"
	DefaultThroughputTarget = "https://speed.cloudflare.com/__down?bytes=10000000"
)

const pingCmd = "ping"

// Doer performs an HTTP request. *http.Client satisfies it; tests substitute a
// fake to avoid real network access.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// ConnOptions configures a connectivity run.
type ConnOptions struct {
	PingTarget       string
	DNSTarget        string
	WebTarget        string
	ThroughputTarget string
	Throughput       bool
	Count            int
	Timeout          time.Duration
}

// withDefaults fills unset fields with package defaults.
func (o ConnOptions) withDefaults() ConnOptions {
	if o.PingTarget == "" {
		o.PingTarget = DefaultPingTarget
	}
	if o.DNSTarget == "" {
		o.DNSTarget = DefaultDNSTarget
	}
	if o.WebTarget == "" {
		o.WebTarget = DefaultWebTarget
	}
	if o.ThroughputTarget == "" {
		o.ThroughputTarget = DefaultThroughputTarget
	}
	if o.Count <= 0 {
		o.Count = 4
	}
	if o.Timeout <= 0 {
		o.Timeout = 5 * time.Second
	}
	return o
}

// Connectivity runs basic reachability tests: gateway/internet/DNS ping, an
// HTTPS GET, and an optional throughput download.
type Connectivity struct {
	Runner run.Runner
	HTTP   Doer
}

// NewConnectivity returns a Connectivity collector. A nil doer defaults to a
// plain http.Client.
func NewConnectivity(r run.Runner, doer Doer) *Connectivity {
	if doer == nil {
		doer = &http.Client{}
	}
	return &Connectivity{Runner: r, HTTP: doer}
}

// Run executes the test suite and returns one result per test.
func (c *Connectivity) Run(ctx context.Context, opt ConnOptions) []model.TestResult {
	opt = opt.withDefaults()

	results := []model.TestResult{
		c.gatewayPing(ctx, opt),
		c.pingTest(ctx, "internet_ping", opt.PingTarget, opt),
		c.pingTest(ctx, "dns_ping", opt.DNSTarget, opt),
		c.webGet(ctx, opt),
	}
	if opt.Throughput {
		results = append(results, c.throughput(ctx, opt))
	}
	return results
}

// gateway resolves the next-hop gateway for target via `ip route get`.
func (c *Connectivity) gateway(ctx context.Context, target string) string {
	out, _, err := c.Runner.Run(ctx, ipCmd, "--json", "route", "get", target)
	if err != nil {
		return ""
	}
	routes, perr := parse.IPRoutes(out)
	if perr != nil || len(routes) == 0 {
		return ""
	}
	return routes[0].Gateway
}

// gatewayPing pings the default gateway, skipping when the route is direct.
func (c *Connectivity) gatewayPing(ctx context.Context, opt ConnOptions) model.TestResult {
	gw := c.gateway(ctx, opt.PingTarget)
	if gw == "" {
		return model.TestResult{
			Name:    "gateway_ping",
			Status:  model.StatusSkip,
			Details: "no gateway for target (direct or link-scoped route)",
		}
	}
	return c.pingTest(ctx, "gateway_ping", gw, opt)
}

// pingTest pings target and reports pass/fail with average latency.
func (c *Connectivity) pingTest(ctx context.Context, name, target string, opt ConnOptions) model.TestResult {
	res := model.TestResult{Name: name, Target: target}
	if target == "" {
		res.Status = model.StatusSkip
		res.Details = "no target"
		return res
	}

	deadline := int(opt.Timeout.Seconds())
	if deadline <= 0 {
		deadline = 5
	}
	out, _, err := c.Runner.Run(ctx, pingCmd,
		"-q", "-c", strconv.Itoa(opt.Count), "-w", strconv.Itoa(deadline), target)

	if avg, ok := parse.PingAvgMillis(out); ok {
		res.Latency = avg + " ms"
	}
	switch {
	case err == nil:
		res.Status = model.StatusPass
	case run.IsNotFound(err):
		res.Status = model.StatusSkip
		res.Error = "ping not found"
	default:
		res.Status = model.StatusFail
		res.Error = "no reply"
	}
	return res
}

// webGet performs an HTTPS GET and passes on a 2xx/3xx status.
func (c *Connectivity) webGet(ctx context.Context, opt ConnOptions) model.TestResult {
	res := model.TestResult{Name: "web_get", Target: opt.WebTarget}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, opt.WebTarget, nil)
	if err != nil {
		res.Status = model.StatusFail
		res.Error = err.Error()
		return res
	}

	start := time.Now()
	resp, err := c.HTTP.Do(req)
	if err != nil {
		res.Status = model.StatusFail
		res.Error = err.Error()
		return res
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<20))
	res.Latency = time.Since(start).Round(time.Millisecond).String()
	res.Details = "HTTP " + resp.Status

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		res.Status = model.StatusPass
	} else {
		res.Status = model.StatusFail
		res.Error = resp.Status
	}
	return res
}

// throughput downloads the target and reports an approximate download rate.
func (c *Connectivity) throughput(ctx context.Context, opt ConnOptions) model.TestResult {
	res := model.TestResult{Name: "throughput_download", Target: opt.ThroughputTarget}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, opt.ThroughputTarget, nil)
	if err != nil {
		res.Status = model.StatusFail
		res.Error = err.Error()
		return res
	}

	start := time.Now()
	resp, err := c.HTTP.Do(req)
	if err != nil {
		res.Status = model.StatusFail
		res.Error = err.Error()
		return res
	}
	defer resp.Body.Close()
	n, err := io.Copy(io.Discard, resp.Body)
	elapsed := time.Since(start)
	if err != nil {
		res.Status = model.StatusFail
		res.Error = err.Error()
		return res
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		res.Status = model.StatusFail
		res.Error = resp.Status
		return res
	}

	secs := elapsed.Seconds()
	mbps := 0.0
	if secs > 0 {
		mbps = float64(n) * 8 / 1e6 / secs
	}
	res.Status = model.StatusPass
	res.Latency = elapsed.Round(time.Millisecond).String()
	res.Details = fmt.Sprintf("%.1f Mbps (%d bytes)", mbps, n)
	return res
}
