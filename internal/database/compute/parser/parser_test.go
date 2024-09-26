package parser

import (
	"reflect"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	p := New()

	tests := []struct {
		name     string
		input    string
		expected *Command
		wantErr  bool
	}{
		{
			name:  "Valid SET command",
			input: "SET key1 value1",
			expected: &Command{
				Action: "SET",
				Args:   []string{"key1", "value1"},
			},
			wantErr: false,
		},
		{
			name:  "Valid GET command",
			input: "GET key1",
			expected: &Command{
				Action: "GET",
				Args:   []string{"key1"},
			},
			wantErr: false,
		},
		{
			name:  "Valid DEL command",
			input: "DEL key1",
			expected: &Command{
				Action: "DEL",
				Args:   []string{"key1"},
			},
			wantErr: false,
		},
		{
			name:     "SET with missing argument",
			input:    "SET key1",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "GET with no arguments",
			input:    "GET",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "DEL with no arguments",
			input:    "DEL",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "SET with extra argument",
			input:    "SET key1 value1 value2",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Unknown command",
			input:    "UNKNOWN key1",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: nil,
			wantErr:  true,
		},
		{
			name:  "Command with punctuation",
			input: "SET key-name value_name",
			expected: &Command{
				Action: "SET",
				Args:   []string{"key-name", "value_name"},
			},
			wantErr: false,
		},
		{
			name:  "Command with asterisks",
			input: "DEL user_****",
			expected: &Command{
				Action: "DEL",
				Args:   []string{"user_****"},
			},
			wantErr: false,
		},
		{
			name:     "Argument with spaces",
			input:    "SET key1 value with spaces",
			expected: nil,
			wantErr:  true,
		},
		{
			name:  "Command with extra spaces",
			input: "  SET   key1   value1  ",
			expected: &Command{
				Action: "SET",
				Args:   []string{"key1", "value1"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestParser_FsmParse(t *testing.T) {
	p := New()

	tests := []struct {
		name     string
		input    string
		expected *Command
		wantErr  bool
	}{
		{
			name:     "Command with extra spaces",
			input:    "  SETwe   key1   value1  ",
			expected: &Command{},
			wantErr:  true,
		},
		{
			name:  "Command with extra spaces",
			input: "   SET   key1   value1    ",
			expected: &Command{
				Action: "SET",
				Args:   []string{"key1", "value1"},
			},
			wantErr: false,
		},
		{
			name:  "Command with extra spaces",
			input: "   GET   key1   value1    ",
			expected: &Command{
				Action: "SET",
				Args:   []string{"key1", "value1"},
			},
			wantErr: true,
		},
		{
			name:  "Command with extra spaces",
			input: "   DEL   key1   value1    ",
			expected: &Command{
				Action: "SET",
				Args:   []string{"key1", "value1"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
