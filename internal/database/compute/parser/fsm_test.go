package parser

import (
	"reflect"
	"testing"
)

func TestFSM_Tokenize(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      []string
		expectErr bool
		errMsg    string
	}{
		{
			name:      "Valid Command with No Arguments",
			input:     "HELP",
			want:      []string{"HELP"},
			expectErr: false,
		},
		{
			name:      "Valid Command with One Argument",
			input:     "RUN test",
			want:      []string{"RUN", "test"},
			expectErr: false,
		},
		{
			name:      "Valid Command with Two Arguments",
			input:     "CMD arg1 arg2",
			want:      []string{"CMD", "arg1", "arg2"},
			expectErr: false,
		},
		{
			name:      "Valid Command with Extra Spaces",
			input:     "  CMD   arg1  arg2  ",
			want:      []string{"CMD", "arg1", "arg2"},
			expectErr: false,
		},
		{
			name:      "Invalid Command with Lowercase Letter",
			input:     "command arg1",
			want:      nil,
			expectErr: true,
			errMsg:    "failed transition.Action: invalid character at start: 'c'",
		},
		{
			name:      "Invalid Character in Command",
			input:     "COMmAND arg1",
			want:      nil,
			expectErr: true,
			errMsg:    "failed transition.Action: invalid character in command: 'm'",
		},
		{
			name:      "Invalid Character in Argument",
			input:     "CMD arg@1",
			want:      []string{"CMD", "arg@1"},
			expectErr: false,
			errMsg:    "",
		},
		{
			name:      "Too Many Arguments",
			input:     "CMD arg1 arg2 arg3",
			want:      nil,
			expectErr: true,
			errMsg:    "failed transition.Action: too many arguments",
		},
		{
			name:      "Empty Input",
			input:     "",
			want:      nil,
			expectErr: true,
			errMsg:    "empty command",
		},
		{
			name:      "Only Spaces",
			input:     "    ",
			want:      nil,
			expectErr: true,
			errMsg:    "empty command",
		},
		{
			name:      "Command with Underscore in Argument",
			input:     "CMD arg_1 arg-2",
			want:      []string{"CMD", "arg_1", "arg-2"},
			expectErr: false,
		},
		{
			name:      "Multiple Spaces Between Arguments",
			input:     "CMD    arg1     arg2",
			want:      []string{"CMD", "arg1", "arg2"},
			expectErr: false,
		},
		{
			name:      "Command Only with Spaces",
			input:     "   CMD   ",
			want:      []string{"CMD"},
			expectErr: false,
		},
		{
			name:      "Argument Starts with Digit",
			input:     "CMD 1arg arg2",
			want:      []string{"CMD", "1arg", "arg2"},
			expectErr: false,
		},
		{
			name:      "Invalid Symbol at Start",
			input:     "#CMD arg1",
			want:      nil,
			expectErr: true,
			errMsg:    "failed transition.Action: invalid character at start: '#'",
		},
	}

	for _, tt := range tests {
		tt := tt // Захват переменной для параллельных тестов
		t.Run(tt.name, func(t *testing.T) {
			fsm := NewFSM(tt.input)
			got, err := fsm.Tokenize()

			if tt.expectErr {
				if err == nil {
					t.Errorf("ожидалась ошибка '%s', но ошибки не было", tt.errMsg)
				} else if err.Error() != tt.errMsg {
					t.Errorf("ожидалась ошибка '%s', получена '%s'", tt.errMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("не ожидалась ошибка, но получена: %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("результат не совпадает.\nОжидалось: %v\nПолучено: %v", tt.want, got)
			}
		})
	}
}
