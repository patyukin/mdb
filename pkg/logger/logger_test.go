package logger

import (
	"github.com/patyukin/mdb/internal/config"
	"testing"

	"go.uber.org/zap/zapcore"
)

type MockConfig struct {
	Logger struct {
		Level string
		Mode  string
	}
}

func newMockConfig(level, mode string) *MockConfig {
	cfg := &MockConfig{}
	cfg.Logger.Level = level
	cfg.Logger.Mode = mode

	return cfg
}

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name          string
		config        *MockConfig
		expectError   bool
		expectedLevel zapcore.Level
		expectedMode  string
	}{
		{
			name:          "Development mode with valid level",
			config:        newMockConfig("debug", "devel"),
			expectError:   false,
			expectedLevel: zapcore.DebugLevel,
			expectedMode:  "devel",
		},
		{
			name:          "Production mode with valid level",
			config:        newMockConfig("info", "prod"),
			expectError:   false,
			expectedLevel: zapcore.InfoLevel,
			expectedMode:  "prod",
		},
		{
			name:        "Invalid log level",
			config:      newMockConfig("invalid_level", "devel"),
			expectError: true,
		},
		{
			name:          "Development mode with uppercase level",
			config:        newMockConfig("ERROR", "devel"),
			expectError:   false,
			expectedLevel: zapcore.ErrorLevel,
			expectedMode:  "devel",
		},
		{
			name:          "Production mode with lowercase level",
			config:        newMockConfig("warn", "prod"),
			expectError:   false,
			expectedLevel: zapcore.WarnLevel,
			expectedMode:  "prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := InitLogger((*config.Config)(tt.config))
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			core := logger.Core()
			enabled := core.Enabled(tt.expectedLevel)
			if !enabled {
				t.Errorf("expected log level %v to be enabled", tt.expectedLevel)
			}

			_ = logger.Sync()
		})
	}
}
