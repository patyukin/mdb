package compute

import (
	"errors"
	"github.com/patyukin/mdb/internal/compute/parser"
	"github.com/patyukin/mdb/internal/storage"
	"github.com/patyukin/mdb/internal/storage/engine"
	"log"
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestCompute_ProcessRequest(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult string
		expectedError  string
		errType        error
	}{
		{
			name:           "Parsing error",
			input:          "SELECT users",
			expectedResult: "Executed: SELECT * FROM users",
			expectedError:  "failed c.parser.Parse: unknown command: SELECT",
		},
		{
			name:           "Parsing error",
			input:          "GET users",
			expectedResult: "Executed: SELECT * FROM users",
			expectedError:  "failed c.storage.Execute: failed s.engine.Get, err: 'users' - key not found",
			errType:        engine.ErrNotFound,
		},
		{
			name:           "Parsing error",
			input:          "INVALID COMMAND",
			expectedResult: "",
			expectedError:  "failed c.parser.Parse: unknown command: INVALID",
		},
		{
			name:           "Execution error",
			input:          "GET *",
			expectedResult: "",
			expectedError:  "failed c.storage.Execute: failed s.engine.GetByPattern, err: key not found",
			errType:        engine.ErrNotFound,
		},
		{
			name:           "Execution error",
			input:          "SELECT * error_table",
			expectedResult: "",
			expectedError:  "failed c.parser.Parse: unknown command: SELECT",
		},
		{
			name:           "Empty input",
			input:          "   ",
			expectedResult: "",
			expectedError:  "failed c.parser.Parse: failed f.Parse: empty command",
		},
		{
			name:           "Successful command",
			input:          "SET qwerty 12345",
			expectedResult: "",
			expectedError:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := zap.NewNop()
			p := parser.New()
			e := engine.New()
			st := storage.New(e, l)

			compute := New(p, st, l)

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
				if result != tt.expectedResult {
					t.Errorf("Expected result '%s', got '%s'", tt.expectedResult, result)
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

func TestCompute_Start(t *testing.T) {
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	inputs := []string{
		"GET users\n",
		"INVALID COMMAND\n",
		"GET error_table werwer\n",
		"\n",
		"EXIT\n",
	}

	r, w, _ := os.Pipe()
	os.Stdin = r

	go func() {
		for _, input := range inputs {
			_, err := w.Write([]byte(input))
			if err != nil {
				log.Printf("Failed to write to pipe: %v", err)
			}
		}

		err := w.Close()
		if err != nil {
			log.Printf("Failed to write to pipe: %v", err)
		}
	}()

	l := zap.NewNop()
	p := &parser.Parser{}
	e := &engine.Engine{}
	st := storage.New(e, l)

	compute := New(p, st, l)
	compute.Start()
}
