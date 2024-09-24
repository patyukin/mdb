package storage

import (
	"errors"
	"github.com/patyukin/mdb/internal/database/compute/parser"
	"github.com/patyukin/mdb/internal/database/storage/mocks"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestStorage_Execute(t *testing.T) {
	t.Parallel()
	mockEngine := new(mocks.Engine)
	logger := zap.NewNop()
	storage := New(mockEngine, logger)

	tests := []struct {
		name        string
		command     *parser.Command
		setupMocks  func()
		expected    string
		expectedErr error
	}{
		{
			name: "SET command",
			command: &parser.Command{
				Action: "SET",
				Args:   []string{"key1", "value1"},
			},
			setupMocks: func() {
				mockEngine.On("Set", "key1", "value1").Once()
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			name: "GET command without wildcard",
			command: &parser.Command{
				Action: "GET",
				Args:   []string{"key1"},
			},
			setupMocks: func() {
				mockEngine.On("Get", "key1").Return("value1", nil).Once()
			},
			expected:    "value1",
			expectedErr: nil,
		},
		{
			name: "GET * command",
			command: &parser.Command{
				Action: "GET",
				Args:   []string{"*"},
			},
			setupMocks: func() {
				mockEngine.On("Get", "*").Return("value1", nil).Once()
			},
			expected:    "value1",
			expectedErr: nil,
		},
		{
			name: "GET command with wildcard but no matches",
			command: &parser.Command{
				Action: "GET",
				Args:   []string{"nonexistent*"},
			},
			setupMocks: func() {
				mockEngine.On("Get", "nonexistent*").Return("", errors.New("key not found")).Once()
			},
			expected:    "",
			expectedErr: errors.New("failed s.engine.Get, err: key not found"),
		},
		{
			name: "DEL command without wildcard",
			command: &parser.Command{
				Action: "DEL",
				Args:   []string{"key1"},
			},
			setupMocks: func() {
				mockEngine.On("Delete", "key1").Return(nil).Once()
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			name: "DEL * command",
			command: &parser.Command{
				Action: "DEL",
				Args:   []string{"*"},
			},
			setupMocks: func() {
				mockEngine.On("Delete", "*").Return(nil).Once()
			},
			expected:    "",
			expectedErr: nil,
		},
		{
			name: "DEL * command with error",
			command: &parser.Command{
				Action: "DEL",
				Args:   []string{"*"},
			},
			setupMocks: func() {
				mockEngine.On("Delete", "*").Return(errors.New("del by pattern failed")).Once()
			},
			expected:    "",
			expectedErr: errors.New("failed s.engine.Delete, err: del by pattern failed"),
		},
		{
			name: "Unknown command",
			command: &parser.Command{
				Action: "UNKNOWN",
				Args:   []string{"arg1"},
			},
			setupMocks:  func() {},
			expected:    "",
			expectedErr: errors.New("failed command.Validate, unknown command: UNKNOWN"),
		},
		{
			name: "SET command with insufficient arguments",
			command: &parser.Command{
				Action: "SET",
				Args:   []string{"key1"},
			},
			setupMocks:  func() {},
			expected:    "",
			expectedErr: errors.New("failed command.Validate, 2 arguments required for SET command"),
		},
		{
			name: "GET command with insufficient arguments",
			command: &parser.Command{
				Action: "GET",
				Args:   []string{},
			},
			setupMocks:  func() {},
			expected:    "",
			expectedErr: errors.New("failed command.Validate, command GET requires 1 argument"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			result, err := storage.Execute(tt.command)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestStorage_Execute_EmptyArgs(t *testing.T) {
	mockEngine := new(mocks.Engine)
	logger := zap.NewNop()
	storage := New(mockEngine, logger)

	command := &parser.Command{
		Action: "SET",
		Args:   []string{},
	}

	result, err := storage.Execute(command)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed command.Validate, 2 arguments required for SET command")

	command = &parser.Command{
		Action: "GET",
		Args:   []string{},
	}

	result, err = storage.Execute(command)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed command.Validate, command GET requires 1 argument")

	command = &parser.Command{
		Action: "DEL",
		Args:   []string{},
	}

	result, err = storage.Execute(command)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed command.Validate, command DEL requires 1 argument")
}

func TestStorage_Execute_ExtraArgs(t *testing.T) {
	mockEngine := new(mocks.Engine)
	logger := zap.NewNop()
	storage := New(mockEngine, logger)

	command := &parser.Command{
		Action: "SET",
		Args:   []string{"key1", "value1", "extra"},
	}

	result, err := storage.Execute(command)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed command.Validate, 2 arguments required for SET command")

	command = &parser.Command{
		Action: "GET",
		Args:   []string{"key1", "extra"},
	}

	result, err = storage.Execute(command)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed command.Validate, command GET requires 1 argument")

	command = &parser.Command{
		Action: "DEL",
		Args:   []string{"key1", "extra"},
	}

	result, err = storage.Execute(command)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed command.Validate, command DEL requires 1 argument")
}

func TestStorage_Execute_InvalidArgs(t *testing.T) {
	mockEngine := new(mocks.Engine)
	logger := zap.NewNop()
	storage := New(mockEngine, logger)

	command := &parser.Command{
		Action: "GET",
		Args:   []string{"key@!"},
	}

	mockEngine.On("Get", "key@!").Return("", errors.New("invalid key")).Once()

	result, err := storage.Execute(command)
	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.EqualError(t, err, "failed s.engine.Get, err: invalid key")
}

func TestStorage_Execute_Set_Get(t *testing.T) {
	mockEngine := new(mocks.Engine)
	logger := zap.NewNop()
	storage := New(mockEngine, logger)

	setCommand := &parser.Command{
		Action: "SET",
		Args:   []string{"key1", "value1"},
	}

	mockEngine.On("Set", "key1", "value1").Once()

	result, err := storage.Execute(setCommand)
	assert.NoError(t, err)
	assert.Equal(t, "", result)

	getCommand := &parser.Command{
		Action: "GET",
		Args:   []string{"key1"},
	}

	mockEngine.On("Get", "key1").Return("value1", nil).Once()

	result, err = storage.Execute(getCommand)
	assert.NoError(t, err)
	assert.Equal(t, "value1", result)
}

func TestStorage_Execute_Logger(t *testing.T) {
	logger := zap.NewExample()
	defer func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			log.Printf("failed logger.Sync, err: %v", err)
		}
	}(logger)

	mockEngine := new(mocks.Engine)
	storage := New(mockEngine, logger)

	command := &parser.Command{
		Action: "SET",
		Args:   []string{"key1", "value1"},
	}

	mockEngine.On("Set", "key1", "value1").Once()

	result, err := storage.Execute(command)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestStorage_Execute_MultipleCommands(t *testing.T) {
	mockEngine := new(mocks.Engine)
	logger := zap.NewNop()
	storage := New(mockEngine, logger)

	setCommand := &parser.Command{
		Action: "SET",
		Args:   []string{"key1", "value1"},
	}

	mockEngine.On("Set", "key1", "value1").Once()

	result, err := storage.Execute(setCommand)
	assert.NoError(t, err)
	assert.Equal(t, "", result)

	getCommand := &parser.Command{
		Action: "GET",
		Args:   []string{"key1"},
	}

	mockEngine.On("Get", "key1").Return("value1", nil).Once()

	result, err = storage.Execute(getCommand)
	assert.NoError(t, err)
	assert.Equal(t, "value1", result)

	delCommand := &parser.Command{
		Action: "DEL",
		Args:   []string{"key1"},
	}

	mockEngine.On("Delete", "key1").Return(nil).Once()

	result, err = storage.Execute(delCommand)
	assert.NoError(t, err)
	assert.Equal(t, "", result)
}
