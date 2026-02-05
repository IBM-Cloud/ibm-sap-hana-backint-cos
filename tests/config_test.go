package tests

import (
	"testing"
)

// Test functions

func TestCheckConfig(t *testing.T) {
	configTests := setupConfigTests(t)
	executablePath := getExecutablePath()

	for _, cfgTest := range configTests {
		t.Run(cfgTest.name, func(t *testing.T) {
			output, err := runCommand(t, executablePath, "-p", cfgTest.cfgPath, "-check")

			if cfgTest.shouldSucceed {
				if err != nil {
					t.Fatal("hdbbackint -check should succeed but failed")
				} else {
					t.Log("Config check passed")
				}
			} else {
				if err == nil {
					t.Fatal("hdbbackint -check should fail but succeeded")
				} else {
					if errMsgOk(string(output), cfgTest.msgToCheck) {
						t.Log("Config check passed")
					} else {
						t.Fatalf("Wrong error message(s): \nOutput: %s", output)
					}
				}
			}
		})
	}
}
