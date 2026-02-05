package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	binaryPath        = "../bin/hdbbackint"
	templateDir       = "testdata/configs"
	apiKeyPlaceholder = "{{APIKEYFILE}}"
)

const (
	cfgSuccessFile            = "cfg_success.tpl"
	cfgWrongSectionFile       = "cfg_wrong_section_name.tpl"
	cfgwrongRangeFile         = "cfg_wrong_range.tpl"
	cfgWrongObjectTagsFile    = "cfg_wrong_object_tags.tpl"
	cfgUnknownParmFile        = "cfg_unknown_parm.tpl"
	cfgParmInWrongSectionFile = "cfg_parm_in_wrong_section.tpl"
	cfgObjLockWrongPeriodFile = "cfg_obj_lock_wrong_period.tpl"
	cfgObjMissingPeriodFile   = "cfg_obj_lock_missing_period.tpl"
	cfgMissingMandParmFile    = "cfg_missing_mand_parm.tpl"
	cfgMixedErrorFile         = "cfg_mixed_error.tpl"
)

type ConfigTest struct {
	name          string
	cfgPath       string
	shouldSucceed bool
	msgToCheck    string
}

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
	apikey := os.Getenv("APIKEY")
	if apikey == "" {
		t.Fatal("APIKEY not set.")
	}
	if err := os.WriteFile(apiKeyFile, []byte(apikey), 0644); err != nil {
		t.Fatalf("failed to create API key file: %v", err)
	}
	return apiKeyFile
}

// prepareConfigFile reads a config template, replaces placeholders, and writes it to a temp location
func prepareConfigFile(t *testing.T, tmpDir, apiKeyFile string, templateFile string) string {
	t.Helper()
	templatePath := filepath.Join(templateDir, templateFile)
	cfgFile := strings.ReplaceAll(templateFile, "tpl", "cfg")
	cfgBytes, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read config template %s: %v", templatePath, err)
	}

	cfgContent := strings.ReplaceAll(string(cfgBytes), apiKeyPlaceholder, apiKeyFile)

	cfgPath := filepath.Join(tmpDir, cfgFile)
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	return cfgPath
}

// getExecutablePath returns the path to the hdbbackint binary
func getExecutablePath() string {
	return binaryPath
}

func setupConfigTests(t *testing.T) []ConfigTest {
	t.Helper()

	tmpDir := setupTempDir(t)
	apiKeyFile := createAPIKeyFile(t, tmpDir)

	var configTests []ConfigTest

	configTests = append(configTests, ConfigTest{
		name:          "Valid Configuration",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgSuccessFile),
		shouldSucceed: true,
		msgToCheck:    "",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Invalid section",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgWrongSectionFile),
		shouldSucceed: false,
		msgToCheck: "ERROR: You specified the section 'cloud'," +
			" but it is not part of the hdbbackint configuration.",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Invalid range",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgwrongRangeFile),
		shouldSucceed: false,
		msgToCheck: "ERROR: 'max_concurrency':" +
			" the value '1234' you specified is invalid.",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Invalid object tags",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgWrongObjectTagsFile),
		shouldSucceed: false,
		msgToCheck: "ERROR: 'object_tags': the value" +
			" 'key1,val1,key2=val2' you specified is invalid.",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Unknown parameter",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgUnknownParmFile),
		shouldSucceed: false,
		msgToCheck: "ERROR: You specified 'unknown_key'" +
			" in section 'objects', but the key is unknown." +
			" The value of 'unknown_key' will be ignored.",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Parameter in wrong section",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgParmInWrongSectionFile),
		shouldSucceed: false,
		msgToCheck: "ERROR: You specified 'object_tags' in section" +
			" 'cloud_storage', but key belongs to section objects." +
			" The value of 'object_tags' will be ignored.",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Wrong object lock retention period",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgObjLockWrongPeriodFile),
		shouldSucceed: false,
		msgToCheck: "The value you specified for 'object_lock_retention_period'" +
			" does not have the correct format.",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Missing object lock retention period",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgObjMissingPeriodFile),
		shouldSucceed: false,
		msgToCheck: "ERROR: You specified 'object_lock_retention_mode = cmp'," +
			" but no 'object_lock_retention_period' is specified.",
	})

	configTests = append(configTests, ConfigTest{
		name:          "Missing mandatory parameter",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgMissingMandParmFile),
		shouldSucceed: false,
		msgToCheck:    "ERROR: You did not specify a value for the mandatory parameter",
	})

	configTests = append(configTests, ConfigTest{
		name:          "More than one error",
		cfgPath:       prepareConfigFile(t, tmpDir, apiKeyFile, cfgMixedErrorFile),
		shouldSucceed: false,
		msgToCheck: "ERROR: 'object_tags': the value" +
			" 'key1,val1,key2=val2' you specified is invalid.",
	})

	return configTests
}

// runCommand executes a command and returns output and error
func runCommand(t *testing.T, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func errMsgOk(output string, message string) bool {
	return strings.Contains(output, message)
}
