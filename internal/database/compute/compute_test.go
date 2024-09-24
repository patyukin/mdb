package compute

import (
	"errors"
	"github.com/patyukin/mdb/internal/database/compute/parser"
	"testing"

	"go.uber.org/zap"
)

func TestCompute_ProcessRequest(t *testing.T) {
	tests := []struct {
		name                 string
		input                string
		expectedResultAction string
		expectedError        string
		errType              error
	}{
		{
			name:                 "Parsing error",
			input:                "SELECT users",
			expectedResultAction: "Executed: SELECT * FROM users",
			expectedError:        "failed c.parser.Parse: failed cmd.Validate: unknown command: SELECT",
		},
		{
			name:                 "Parsing error",
			input:                "GET users",
			expectedResultAction: "GET",
			expectedError:        "",
		},
		{
			name:                 "Parsing error",
			input:                "INVALID COMMAND",
			expectedResultAction: "",
			expectedError:        "failed c.parser.Parse: failed cmd.Validate: unknown command: INVALID",
		},
		{
			name:                 "Execution error",
			input:                "GET *",
			expectedResultAction: "GET",
			expectedError:        "",
		},
		{
			name:                 "Execution error",
			input:                "SELECT * error_table",
			expectedResultAction: "",
			expectedError:        "failed c.parser.Parse: failed cmd.Validate: unknown command: SELECT",
		},
		{
			name:                 "Empty input",
			input:                "   ",
			expectedResultAction: "",
			expectedError:        "failed c.parser.Parse: failed f.Parse: empty command",
		},
		{
			name:                 "Successful command",
			input:                "SET qwerty 12345",
			expectedResultAction: "SET",
			expectedError:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := zap.NewNop()
			p := parser.New()

			compute := New(p, l)

			result, err := compute.ProcessRequest(tt.input)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error '%s', got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("Expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got '%s'", err.Error())
				}
				if result.Action != tt.expectedResultAction {
					t.Errorf("Expected result '%s', got '%s'", tt.expectedResultAction, result.Action)
				}
			}

			if tt.errType != nil {
				if err == nil || !errors.Is(err, tt.errType) {
					t.Errorf("Expected error '%s', got '%s'", tt.errType, err)
				}
			}
		})
	}
}
