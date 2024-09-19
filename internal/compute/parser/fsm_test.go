package parser

import (
	"reflect"
	"testing"
)

func TestFSM_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid SET command",
			input:    "SET key value",
			expected: []string{"SET", "key", "value"},
			wantErr:  false,
		},
		{
			name:     "Valid GET command",
			input:    "GET key",
			expected: []string{"GET", "key"},
			wantErr:  false,
		},
		{
			name:     "Valid DEL command",
			input:    "DEL key",
			expected: []string{"DEL", "key"},
			wantErr:  false,
		},
		{
			name:     "Command with extra spaces",
			input:    "  SET   key   value  ",
			expected: []string{"SET", "key", "value"},
			wantErr:  false,
		},
		{
			name:    "Empty input",
			input:   "",
			wantErr: true,
			errMsg:  "empty command",
		},
		{
			name:     "Invalid character in argument",
			input:    "SET key val@ue",
			wantErr:  false,
			expected: []string{"SET", "key", "val@ue"},
		},
		{
			name:     "SET command missing value",
			input:    "SET key",
			wantErr:  false,
			expected: []string{"SET", "key"},
			errMsg:   "ожидался аргумент",
		},
		{
			name:     "GET command missing key",
			input:    "GET",
			wantErr:  false,
			expected: []string{"GET"},
			errMsg:   "ожидался аргумент",
		},
		{
			name:     "SET command with extra arguments",
			input:    "SET key value extra",
			wantErr:  false,
			expected: []string{"SET", "key", "value", "extra"},
		},
		{
			name:     "Argument with allowed special characters",
			input:    "SET /path/to/key value_with-symbols",
			expected: []string{"SET", "/path/to/key", "value_with-symbols"},
			wantErr:  false,
		},
		{
			name:     "Argument containing spaces",
			input:    "SET key value with spaces",
			wantErr:  false,
			expected: []string{"SET", "key", "value", "with", "spaces"},
			errMsg:   "ожидался аргумент",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFSM(tt.input)
			result, err := f.Parse()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				} else if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected result %v, got %v", tt.expected, result)
				}
			}
		})
	}
}
