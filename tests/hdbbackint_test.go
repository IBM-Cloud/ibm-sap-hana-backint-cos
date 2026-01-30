package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func getTmpDir(t *testing.T) string {
	tmpDir := t.TempDir()
	return tmpDir
}

func createApiKeyFile(t *testing.T, tmpDir string) string {
	apiKeyFile := filepath.Join(tmpDir, "apikey")
	if err := os.WriteFile(apiKeyFile, []byte("dummy"), 0644); err != nil {
		t.Fatalf("cannot create catalog file: %v", err)
	}
	return apiKeyFile
}

func getConfigFile(
	t *testing.T,
	tmpDir string,
	templateFile string,
) string {
	apiKeyFile := createApiKeyFile(t, tmpDir)
	cfgBytes, err := os.ReadFile(templateFile)
	if err != nil {
		t.Fatalf("cannot read config template: %v", err)
	}

	cfgContent := strings.ReplaceAll(
		string(cfgBytes),
		"{{APIKEYFILE}}",
		apiKeyFile,
	)

	cfgPath := filepath.Join(tmpDir, "hdbbackint.cfg")
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatalf("cannot write config: %v", err)
	}
	return cfgPath
}

func getExecutable(t *testing.T) string {
	return filepath.Join("..", "bin", "hdbbackint")
}

func TestVersion(t *testing.T) {
	cmd := exec.Command(getExecutable(t), "-v")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hdbbackint -v failed: %v\nOutput: %s", err, output)
	}

	if len(output) == 0 {
		t.Fatal("expected version output, got empty output")
	}
}

func TestCheckConfigSuccess(t *testing.T) {
	templateFile := "../testdata/test_success.cfg"
	cfgPath := getConfigFile(t, getTmpDir(t), templateFile)
	cmd := exec.Command(
		getExecutable(t),
		"-p", cfgPath,
		"-check",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("hdbbackint -check failed: %v\nOutput: %s", err, output)
	}
}

func TestCheckConfigFail(t *testing.T) {
	templateFile := "../testdata/test_fail.cfg"
	cfgPath := getConfigFile(t, getTmpDir(t), templateFile)

	cmd := exec.Command(
		getExecutable(t),
		"-p", cfgPath,
		"-check",
	)

	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("hdbbackint -check failed: %v\nOutput: %s", err, output)
	}
}
