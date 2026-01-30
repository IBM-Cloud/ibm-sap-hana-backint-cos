package tests

import (
	"os/exec"
	"testing"
)

func TestVersion(t *testing.T) {
	cmd := exec.Command("./hdbbackint", "-v")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hdbbackint -v failed: %v\nOutput: %s", err, output)
	}

	if len(output) == 0 {
		t.Fatal("expected version output, got empty output")
	}
}

func TestCheckConfig(t *testing.T) {
	cfgPath := "../testdata/hdbbackint.cfg"

	cmd := exec.Command(
		"./hdbbackint",
		"-p", cfgPath,
		"-check",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hdbbackint -check failed: %v\nOutput: %s", err, output)
	}
}
