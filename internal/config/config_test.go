package config

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func createTempYAML(t *testing.T, content string) (string, func()) {
	tmpDir := os.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "config_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err = tmpFile.Write([]byte(content)); err != nil {
		err = tmpFile.Close()
		if err != nil {
			return "", nil
		}

		err = os.Remove(tmpFile.Name())
		if err != nil {
			return "", nil
		}

		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err = tmpFile.Close(); err != nil {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return "", nil
		}

		t.Fatalf("Failed to close temp file: %v", err)
	}

	cleanup := func() {
		err = os.Remove(tmpFile.Name())
		if err != nil {
			return
		}
	}

	return tmpFile.Name(), cleanup
}

func TestLoadConfig_Success(t *testing.T) {
	yamlContent := `
logger:
  level: "info"
  mode: "prod"
`

	filePath, cleanup := createTempYAML(t, yamlContent)
	defer cleanup()

	if err := os.Setenv("YAML_CONFIG_FILE_PATH", filePath); err != nil {
		return
	}

	defer func() {
		if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
			log.Printf("failed os.Unsetenv, err: %v", err)
		}
	}()

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Logger.Level != "info" {
		t.Errorf("Expected logger level 'info', got '%s'", config.Logger.Level)
	}

	if config.Logger.Mode != "prod" {
		t.Errorf("Expected logger mode 'prod', got '%s'", config.Logger.Mode)
	}
}

func TestLoadConfig_NoEnvVar(t *testing.T) {
	if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
		return
	}

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected error when YAML_CONFIG_FILE_PATH is not set, got nil")
	}

	expectedErrMsg := "yaml config file path is not set"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	nonExistentFile := filepath.Join(os.TempDir(), "non_existent_config.yaml")

	if err := os.Setenv("YAML_CONFIG_FILE_PATH", nonExistentFile); err != nil {
		return
	}
	defer func() {
		if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
			log.Printf("failed os.Unsetenv, err: %v", err)
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected error when config file does not exist, got nil")
	}

	expectedErrPrefix := "unable to open config file"
	if len(err.Error()) < len(expectedErrPrefix) || err.Error()[:len(expectedErrPrefix)] != expectedErrPrefix {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedErrPrefix, err.Error())
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	yamlContent := `
logger:
  level: "info
  mode: "prod"
`

	filePath, cleanup := createTempYAML(t, yamlContent)
	defer cleanup()

	if err := os.Setenv("YAML_CONFIG_FILE_PATH", filePath); err != nil {
		return
	}

	defer func() {
		if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
			log.Printf("failed os.Unsetenv, err: %v", err)
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected error due to invalid YAML, got nil")
	}

	expectedErrPrefix := "unable to decode config file"
	if len(err.Error()) < len(expectedErrPrefix) || err.Error()[:len(expectedErrPrefix)] != expectedErrPrefix {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedErrPrefix, err.Error())
	}
}

func TestLoadConfig_InvalidValidation_MissingFields(t *testing.T) {
	yamlContent := `
logger:
  level: "info"
`

	filePath, cleanup := createTempYAML(t, yamlContent)
	defer cleanup()

	if err := os.Setenv("YAML_CONFIG_FILE_PATH", filePath); err != nil {
		return
	}

	defer func() {
		if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
			log.Printf("failed os.Unsetenv, err: %v", err)
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected validation error due to missing fields, got nil")
	}

	expectedErrPrefix := "config validation failed"
	if len(err.Error()) < len(expectedErrPrefix) || err.Error()[:len(expectedErrPrefix)] != expectedErrPrefix {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedErrPrefix, err.Error())
	}
}

func TestLoadConfig_InvalidValidation_InvalidValues(t *testing.T) {
	yamlContent := `
logger:
  level: "verbose"
  mode: "staging"
`

	filePath, cleanup := createTempYAML(t, yamlContent)
	defer cleanup()

	if err := os.Setenv("YAML_CONFIG_FILE_PATH", filePath); err != nil {
		return
	}

	defer func() {
		if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
			log.Printf("failed os.Unsetenv, err: %v", err)
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected validation error due to invalid enum values, got nil")
	}

	expectedErrPrefix := "config validation failed"
	if len(err.Error()) < len(expectedErrPrefix) || err.Error()[:len(expectedErrPrefix)] != expectedErrPrefix {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedErrPrefix, err.Error())
	}
}

func TestLoadConfig_InvalidValidation_EmptyLevel(t *testing.T) {
	yamlContent := `
logger:
  level: ""
  mode: "prod"
`

	filePath, cleanup := createTempYAML(t, yamlContent)
	defer cleanup()

	if err := os.Setenv("YAML_CONFIG_FILE_PATH", filePath); err != nil {
		return
	}

	defer func() {
		if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
			log.Printf("failed os.Unsetenv, err: %v", err)
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected validation error due to empty logger.level, got nil")
	}

	expectedErrPrefix := "config validation failed"
	if len(err.Error()) < len(expectedErrPrefix) || err.Error()[:len(expectedErrPrefix)] != expectedErrPrefix {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedErrPrefix, err.Error())
	}
}

func TestLoadConfig_InvalidValidation_EmptyMode(t *testing.T) {
	yamlContent := `
logger:
  level: "info"
  mode: ""
`

	filePath, cleanup := createTempYAML(t, yamlContent)
	defer cleanup()

	if err := os.Setenv("YAML_CONFIG_FILE_PATH", filePath); err != nil {
		return
	}

	defer func() {
		if err := os.Unsetenv("YAML_CONFIG_FILE_PATH"); err != nil {
			log.Printf("failed os.Unsetenv, err: %v", err)
		}
	}()

	_, err := LoadConfig()
	if err == nil {
		t.Fatalf("Expected validation error due to empty logger.mode, got nil")
	}

	expectedErrPrefix := "config validation failed"
	if len(err.Error()) < len(expectedErrPrefix) || err.Error()[:len(expectedErrPrefix)] != expectedErrPrefix {
		t.Errorf("Expected error message to start with '%s', got '%s'", expectedErrPrefix, err.Error())
	}
}
