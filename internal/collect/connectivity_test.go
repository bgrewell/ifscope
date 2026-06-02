package collect

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/testutil"
)

// errExit models a non-zero process exit (as ping returns on packet loss).
var errExit = errors.New("exit status 1")

// fakeDoer returns a canned HTTP response (or error) for every request.
type fakeDoer struct {
	resp *http.Response
	err  error
}

func (f fakeDoer) Do(*http.Request) (*http.Response, error) {
	return f.resp, f.err
}

func httpResp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func resultsByName(rs []model.TestResult) map[string]model.TestResult {
	m := make(map[string]model.TestResult, len(rs))
	for _, r := range rs {
		m[r.Name] = r
	}
	return m
}

func TestConnectivityAllPass(t *testing.T) {
	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: `[{"dst":"1.1.1.1","gateway":"172.26.1.1","dev":"eth0"}]`},
		"ip", "--json", "route", "get", "1.1.1.1")
	okPing := run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/ok.txt"))}
	fake.SetResult(okPing, "ping", "-q", "-c", "4", "-w", "5", "172.26.1.1")
	fake.SetResult(okPing, "ping", "-q", "-c", "4", "-w", "5", "1.1.1.1")
	fake.SetResult(okPing, "ping", "-q", "-c", "4", "-w", "5", "www.cloudflare.com")

	doer := fakeDoer{resp: httpResp(200, "ok")}
	results := NewConnectivity(fake, doer).Run(context.Background(), ConnOptions{})
	by := resultsByName(results)

	for _, name := range []string{"gateway_ping", "internet_ping", "dns_ping", "web_get"} {
		if by[name].Status != model.StatusPass {
			t.Errorf("%s status = %q, want pass (%+v)", name, by[name].Status, by[name])
		}
	}
	if by["gateway_ping"].Target != "172.26.1.1" {
		t.Errorf("gateway target = %q", by["gateway_ping"].Target)
	}
	if by["internet_ping"].Latency != "6.154 ms" {
		t.Errorf("latency = %q, want 6.154 ms", by["internet_ping"].Latency)
	}
}

func TestConnectivityGatewaySkippedWhenDirect(t *testing.T) {
	fake := run.NewFake()
	// route get returns no gateway (link-scoped/direct).
	fake.SetResult(run.FakeResult{Stdout: `[{"dst":"1.1.1.1","dev":"eth0","prefsrc":"1.1.1.2"}]`},
		"ip", "--json", "route", "get", "1.1.1.1")
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/ok.txt"))},
		"ping", "-q", "-c", "4", "-w", "5", "1.1.1.1")
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/ok.txt"))},
		"ping", "-q", "-c", "4", "-w", "5", "www.cloudflare.com")

	results := NewConnectivity(fake, fakeDoer{resp: httpResp(200, "ok")}).Run(context.Background(), ConnOptions{})
	by := resultsByName(results)
	if by["gateway_ping"].Status != model.StatusSkip {
		t.Errorf("gateway_ping = %q, want skip", by["gateway_ping"].Status)
	}
}

func TestConnectivityPingFails(t *testing.T) {
	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: `[{"dst":"1.1.1.1","gateway":"172.26.1.1"}]`},
		"ip", "--json", "route", "get", "1.1.1.1")
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/ok.txt"))},
		"ping", "-q", "-c", "4", "-w", "5", "172.26.1.1")
	// internet ping returns loss and a non-nil error (exit status).
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/loss.txt")), Err: errExit},
		"ping", "-q", "-c", "4", "-w", "5", "1.1.1.1")
	fake.SetResult(run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/ok.txt"))},
		"ping", "-q", "-c", "4", "-w", "5", "www.cloudflare.com")

	results := NewConnectivity(fake, fakeDoer{resp: httpResp(200, "ok")}).Run(context.Background(), ConnOptions{})
	by := resultsByName(results)
	if by["internet_ping"].Status != model.StatusFail {
		t.Errorf("internet_ping = %q, want fail", by["internet_ping"].Status)
	}
}

func TestConnectivityWebFailOnServerError(t *testing.T) {
	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: `[{"dst":"1.1.1.1","gateway":"172.26.1.1"}]`},
		"ip", "--json", "route", "get", "1.1.1.1")
	ok := run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/ok.txt"))}
	fake.SetResult(ok, "ping", "-q", "-c", "4", "-w", "5", "172.26.1.1")
	fake.SetResult(ok, "ping", "-q", "-c", "4", "-w", "5", "1.1.1.1")
	fake.SetResult(ok, "ping", "-q", "-c", "4", "-w", "5", "www.cloudflare.com")

	results := NewConnectivity(fake, fakeDoer{resp: httpResp(503, "down")}).Run(context.Background(), ConnOptions{})
	if resultsByName(results)["web_get"].Status != model.StatusFail {
		t.Errorf("web_get should fail on 503")
	}
}

func TestThroughputOptIn(t *testing.T) {
	fake := run.NewFake()
	fake.SetResult(run.FakeResult{Stdout: `[{"dst":"1.1.1.1","gateway":"172.26.1.1"}]`},
		"ip", "--json", "route", "get", "1.1.1.1")
	ok := run.FakeResult{Stdout: string(testutil.Fixture(t, "ping/ok.txt"))}
	fake.SetResult(ok, "ping", "-q", "-c", "4", "-w", "5", "172.26.1.1")
	fake.SetResult(ok, "ping", "-q", "-c", "4", "-w", "5", "1.1.1.1")
	fake.SetResult(ok, "ping", "-q", "-c", "4", "-w", "5", "www.cloudflare.com")
	doer := fakeDoer{resp: httpResp(200, strings.Repeat("x", 1024))}

	without := NewConnectivity(fake, doer).Run(context.Background(), ConnOptions{})
	if _, ok := resultsByName(without)["throughput_download"]; ok {
		t.Error("throughput should be opt-in")
	}
	with := NewConnectivity(fake, doer).Run(context.Background(), ConnOptions{Throughput: true})
	if _, ok := resultsByName(with)["throughput_download"]; !ok {
		t.Error("throughput should run when requested")
	}
}
