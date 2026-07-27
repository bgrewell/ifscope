package app

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// Keep the user-facing command inventory synchronized with the Cobra tree.
// This intentionally checks README.md rather than duplicating the command list
// in another source file: adding a command without documenting it should fail.
func TestREADMECommandsMatchCobraTree(t *testing.T) {
	readme, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}

	documented := make(map[string]bool)
	row := regexp.MustCompile("(?m)^\\| `ifscope(?: ([a-z][a-z0-9-]*))?`(?: / `ifscope ([a-z][a-z0-9-]*)`)?")
	for _, match := range row.FindAllStringSubmatch(string(readme), -1) {
		for _, name := range match[1:] {
			if name != "" {
				documented[name] = true
			}
		}
	}

	root := newRootCommand(&Options{})
	for _, command := range root.Commands() {
		switch command.Name() {
		case "completion", "help":
			continue
		}
		if !documented[command.Name()] {
			t.Errorf("README command table does not document %q", command.Name())
		}
	}

	root.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		switch flag.Name {
		case "help", "version":
			return
		}
		if !strings.Contains(string(readme), "--"+flag.Name) {
			t.Errorf("README does not document global flag --%s", flag.Name)
		}
	})
}
