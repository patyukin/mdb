package database

import (
	"errors"
	"github.com/patyukin/mdb/internal/database/compute"
	"github.com/patyukin/mdb/internal/database/compute/parser"
	"github.com/patyukin/mdb/internal/database/mocks"
	"github.com/patyukin/mdb/internal/database/storage"
	"github.com/patyukin/mdb/internal/database/storage/engine"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestHandleQuery_TableDriven(t *testing.T) {
	logger := zap.NewNop()
	mockStorage := new(mocks.Storage)
	mockCompute := new(mocks.Compute)
	db := New(mockCompute, mockStorage, logger)

	tests := []struct {
		name           string
		request        string
		setupMocks     func()
		expectedResult string
		expectError    bool
		errorMessage   string
	}{
		{
			name:    "Successful HandleQuery",
			request: "SET test_key test_value",
			setupMocks: func() {
				cmd := &parser.Command{Action: "SET", Args: []string{"test_key", "test_value"}}
				result := ""
				mockCompute.On("ProcessRequest", "SET test_key test_value").Return(cmd, nil)
				mockStorage.On("Execute", cmd).Return(result, nil)
			},
			expectedResult: "",
			expectError:    false,
		},
		{
			name:    "ProcessRequest Failure",
			request: "test_request",
			setupMocks: func() {
				mockCompute.On("ProcessRequest", "test_request").Return(nil, errors.New("process error"))
			},
			expectedResult: "",
			expectError:    true,
			errorMessage:   "failed d.cmpt.ProcessRequest",
		},
		{
			name:    "Execute Failure",
			request: "test_request",
			setupMocks: func() {
				cmd := &parser.Command{Action: "TestCommand"}
				mockCompute.On("ProcessRequest", "test_request").Return(cmd, nil)
				mockStorage.On("Execute", cmd).Return("", errors.New("execute error"))
			},
			expectedResult: "",
			expectError:    true,
			errorMessage:   "failed c.storage.Execute",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			res, err := db.HandleQuery(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Equal(t, tt.expectedResult, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, res)
			}

			mockCompute.AssertExpectations(t)
			mockStorage.AssertExpectations(t)

			mockCompute.ExpectedCalls = nil
			mockCompute.Calls = nil
			mockStorage.ExpectedCalls = nil
			mockStorage.Calls = nil
		})
	}
}

func TestHandleQuery(t *testing.T) {
	logger := zap.NewNop()

	prsr := parser.New()
	cmpt := compute.New(prsr, logger)
	eng := engine.New()
	strg := storage.New(eng, logger)

	db := New(cmpt, strg, logger)

	tests := []struct {
		name           string
		request        string
		setupStorage   func()
		expectedResult string
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "Успешная обработка запроса",
			request:        "SET test_key test_value",
			setupStorage:   func() {},
			expectedResult: "",
			expectError:    false,
		},
		{
			name:           "Ошибка при обработке запроса",
			request:        "",
			setupStorage:   func() {},
			expectedResult: "",
			expectError:    true,
			errorMessage:   "failed d.cmpt.ProcessRequest",
		},
		{
			name:           "Ошибка при выполнении команды",
			request:        "unknown_command",
			setupStorage:   func() {},
			expectedResult: "",
			expectError:    true,
			errorMessage:   "failed d.cmpt.ProcessRequest: failed c.parser.Parse: failed f.Parse: failed transition.Action: invalid character at start: 'u'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupStorage()

			res, err := db.HandleQuery(tt.request)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Equal(t, tt.expectedResult, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, res)
			}
		})
	}
}
