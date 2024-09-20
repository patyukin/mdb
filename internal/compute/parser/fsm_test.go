package parser

import (
	"testing"
)

// Тестовая функция для stateCommandToken
func TestStateCommandToken(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		initialPos  int
		expectedErr string
		expectedTok []string
		finalPos    int
	}{
		{
			name:        "Valid command with arguments",
			input:       "SET key value",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"SET"},
			finalPos:    3,
		},
		{
			name:        "Valid command with multiple arguments",
			input:       "GET /etc/nginx/config",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"GET"},
			finalPos:    3,
		},
		{
			name:        "Invalid first token lowercase",
			input:       "Set key value",
			initialPos:  0,
			expectedErr: "invalid character in command: 'e'",
			expectedTok: nil,
			finalPos:    1,
		},
		{
			name:        "Invalid first token lowercase",
			input:       "Set",
			initialPos:  0,
			expectedErr: "invalid character in command: 'e'",
			expectedTok: nil,
			finalPos:    1,
		},
		{
			name:        "Invalid character in command",
			input:       "S3T key value",
			initialPos:  0,
			expectedErr: "invalid character in command: '3'",
			expectedTok: nil,
			finalPos:    1,
		},
		{
			name:        "Empty input",
			input:       "",
			initialPos:  0,
			expectedErr: "empty input",
			expectedTok: nil,
			finalPos:    0, // Ошибка на позиции 0
		},
		{
			name:        "Command followed by end of input",
			input:       "DEL",
			initialPos:  0,
			expectedErr: "empty arguments",
			expectedTok: []string{"DEL"},
			finalPos:    3,
		},
		{
			name:        "Command with trailing spaces",
			input:       "   GET   ",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"GET"},
			finalPos:    6,
		},
		{
			name:        "Command with non-English uppercase letters",
			input:       "SET WEATHER_2_PM COLD_MOSCOW_WEATHER",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"SET"},
			finalPos:    3,
		},
		{
			name:        "Command with empty arguments",
			input:       "SET",
			initialPos:  0,
			expectedErr: "empty arguments",
			expectedTok: []string{"SET"},
			finalPos:    3,
		},
		{
			name:        "Empty command with spaces",
			input:       "    ",
			initialPos:  0,
			expectedErr: "empty input",
			expectedTok: nil,
			finalPos:    4,
		},
		{
			name:        "Command with mixed case letters",
			input:       "SeT key value",
			initialPos:  0,
			expectedErr: "invalid character in command: 'e'",
			expectedTok: nil,
			finalPos:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый FSM
			f := NewFSM(tt.input)
			f.position = tt.initialPos

			// Вызываем функцию stateCommandToken
			nextState, err := f.stateCommandToken()

			// Проверяем ожидаемую ошибку
			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("expected error '%s', got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("expected error '%s', got '%s'", tt.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			// Проверяем ожидаемые токены
			if tt.expectedTok != nil {
				if len(f.tokens) != len(tt.expectedTok) {
					t.Errorf("expected tokens %v, got %v", tt.expectedTok, f.tokens)
				} else {
					for i, tok := range tt.expectedTok {
						if f.tokens[i] != tok {
							t.Errorf("expected token '%s', got '%s'", tok, f.tokens[i])
						}
					}
				}
			}

			// Проверяем конечную позицию
			if f.position != tt.finalPos {
				t.Errorf("expected final position %d, got %d", tt.finalPos, f.position)
			}

			// Проверяем следующий состояние
			if tt.expectedErr != "" {
				if nextState != nil {
					t.Errorf("expected next state to be nil on error, got %v", nextState)
				}
			}
		})
	}
}

// Тестовая функция для stateArgumentsToken
func TestStateArgumentsToken(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		initialPos  int
		expectedErr string
		expectedTok []string
		finalPos    int
	}{
		{
			name:        "Single valid argument",
			input:       "value1",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"value1"},
			finalPos:    6,
		},
		{
			name:        "Multiple valid arguments",
			input:       "arg1 arg2 arg3",
			initialPos:  0,
			expectedErr: "too many arguments",
			expectedTok: []string{"arg1", "arg2"},
			finalPos:    11,
		},
		{
			name:        "Arguments with punctuation",
			input:       "user_name /path/to/resource key-value",
			initialPos:  0,
			expectedErr: "too many arguments",
			expectedTok: []string{"user_name", "/path/to/resource"},
			finalPos:    29,
		},
		{
			name:        "Arguments with invalid character",
			input:       "validArg invalid@Arg",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"validArg", "invalid@Arg"},
			finalPos:    20,
		},
		{
			name:        "Arguments with trailing spaces",
			input:       "arg1 arg2   ",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"arg1", "arg2"},
			finalPos:    12,
		},
		{
			name:        "Empty input",
			input:       "",
			initialPos:  0,
			expectedErr: "empty arguments",
			expectedTok: []string{},
			finalPos:    0,
		},
		{
			name:        "Only spaces",
			input:       "     ",
			initialPos:  0,
			expectedErr: "empty arguments",
			expectedTok: []string{},
			finalPos:    5,
		},
		{
			name:        "Arguments with mixed valid and invalid characters",
			input:       "arg1 arg$ arg3",
			initialPos:  0,
			expectedErr: "too many arguments",
			expectedTok: []string{"arg1", "arg$"},
			finalPos:    11,
		},
		{
			name:        "Arguments with slash",
			input:       "/path/to/resource arg2",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"/path/to/resource", "arg2"},
			finalPos:    22,
		},
		{
			name:        "Arguments starting with punctuation",
			input:       "*starArg arg2",
			initialPos:  0,
			expectedErr: "",
			expectedTok: []string{"*starArg", "arg2"},
			finalPos:    13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём новый FSM
			f := &FSM{
				input:    tt.input,
				position: tt.initialPos,
				tokens:   []string{"CMD"},
			}

			// Вызываем функцию stateArgumentsToken
			nextState, err := f.stateArgumentsToken()

			f.tokens = f.tokens[1:]

			// Проверяем ожидаемую ошибку
			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("expected error '%s', got nil", tt.expectedErr)
				} else if err.Error() != tt.expectedErr {
					t.Errorf("expected error '%s', got '%s'", tt.expectedErr, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			// Проверяем ожидаемые токены
			if tt.expectedTok != nil {
				if len(f.tokens) != len(tt.expectedTok) {
					t.Errorf("expected tokens %v, got %v", tt.expectedTok, f.tokens)
				} else {
					for i, tok := range tt.expectedTok {
						if f.tokens[i] != tok {
							t.Errorf("expected token '%s', got '%s'", tok, f.tokens[i])
						}
					}
				}
			}

			// Проверяем конечную позицию
			if f.position != tt.finalPos {
				t.Errorf("expected final position %d, got %d", tt.finalPos, f.position)
			}

			// Проверяем следующий состояние
			if tt.expectedErr == "" && nextState != nil {
				// Если ожидается ошибка, nextState должно быть nil
				t.Errorf("expected next state to be nil on error, got %v", nextState)
			}
		})
	}
}
