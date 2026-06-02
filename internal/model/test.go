package model

// TestStatus is the outcome of a single connectivity test.
type TestStatus string

// Connectivity test outcomes.
const (
	StatusPass    TestStatus = "pass"
	StatusFail    TestStatus = "fail"
	StatusSkip    TestStatus = "skip"
	StatusUnknown TestStatus = "unknown"
)

// TestResult is the outcome of one connectivity test.
type TestResult struct {
	Name    string     `json:"name"`
	Status  TestStatus `json:"status"`
	Target  string     `json:"target,omitempty"`
	Latency string     `json:"latency,omitempty"`
	Details string     `json:"details,omitempty"`
	Error   string     `json:"error,omitempty"`
}
