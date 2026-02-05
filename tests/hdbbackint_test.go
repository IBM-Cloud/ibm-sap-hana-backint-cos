package tests

import (
	"testing"
)

// Test functions

func TestVersion(t *testing.T) {
	output, err := runCommand(t, getExecutablePath(), "-v")
	if err != nil {
		t.Fatalf("hdbbackint -v failed: %v\nOutput: %s", err, output)
	}

	if len(output) == 0 {
		t.Fatal("expected version output, got empty output")
	}

	t.Logf("Version output: %s", string(output))
}
