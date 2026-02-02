package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	binaryPath         = "../bin/hdbbackint"
	validConfigPath    = "testdata/valid_config.cfg"
	invalidConfigPath  = "testdata/invalid_config.cfg"
	apiKeyPlaceholder  = "{{APIKEYFILE}}"
	dummyAPIKeyContent = "dummy"
)

// Helper functions

// setupTempDir creates a temporary directory for test files
func setupTempDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// createAPIKeyFile creates a dummy API key file in the specified directory
func createAPIKeyFile(t *testing.T, tmpDir string) string {
	t.Helper()
	apiKeyFile := filepath.Join(tmpDir, "apikey")
	if err := os.WriteFile(apiKeyFile, []byte(dummyAPIKeyContent), 0644); err != nil {
		t.Fatalf("failed to create API key file: %v", err)
	}
	return apiKeyFile
}

// prepareConfigFile reads a config template, replaces placeholders, and writes it to a temp location
func prepareConfigFile(t *testing.T, tmpDir, templatePath string) string {
	t.Helper()

	apiKeyFile := createAPIKeyFile(t, tmpDir)

	cfgBytes, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read config template %s: %v", templatePath, err)
	}

	cfgContent := strings.ReplaceAll(string(cfgBytes), apiKeyPlaceholder, apiKeyFile)

	cfgPath := filepath.Join(tmpDir, "hdbbackint.cfg")
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	return cfgPath
}

// getExecutablePath returns the path to the hdbbackint binary
func getExecutablePath() string {
	return binaryPath
}

// runCommand executes a command and returns output and error
func runCommand(t *testing.T, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

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

func TestCheckConfigSuccess(t *testing.T) {
	tmpDir := setupTempDir(t)
	cfgPath := prepareConfigFile(t, tmpDir, validConfigPath)

	output, err := runCommand(t, getExecutablePath(), "-p", cfgPath, "-check")
	if err != nil {
		t.Fatalf("hdbbackint -check should succeed but failed: %v\nOutput: %s", err, output)
	}

	t.Logf("Config check passed: %s", string(output))
}

func TestCheckConfigFail(t *testing.T) {
	tmpDir := setupTempDir(t)
	cfgPath := prepareConfigFile(t, tmpDir, invalidConfigPath)

	output, err := runCommand(t, getExecutablePath(), "-p", cfgPath, "-check")
	if err == nil {
		t.Fatalf("hdbbackint -check should fail but succeeded\nOutput: %s", output)
	}

	t.Logf("Config check failed as expected: %s", string(output))
}

// Table-driven test for config validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name          string
		configPath    string
		shouldSucceed bool
		description   string
	}{
		{
			name:          "ValidConfig",
			configPath:    validConfigPath,
			shouldSucceed: true,
			description:   "Valid configuration with all required fields",
		},
		{
			name:          "InvalidConfig",
			configPath:    invalidConfigPath,
			shouldSucceed: false,
			description:   "Invalid configuration missing required region field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := setupTempDir(t)
			cfgPath := prepareConfigFile(t, tmpDir, tt.configPath)

			output, err := runCommand(t, getExecutablePath(), "-p", cfgPath, "-check")

			if tt.shouldSucceed && err != nil {
				t.Errorf("%s: expected success but got error: %v\nOutput: %s",
					tt.description, err, output)
			}

			if !tt.shouldSucceed && err == nil {
				t.Errorf("%s: expected failure but got success\nOutput: %s",
					tt.description, output)
			}
		})
	}
}
