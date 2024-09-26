package storage

import (
	"fmt"
	"github.com/patyukin/mdb/internal/database/compute/parser"
	"go.uber.org/zap"
)

const (
	SET    = "SET"
	GET    = "GET"
	DELETE = "DEL"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=Engine --output ./mocks
type Engine interface {
	Set(key string, value string)
	Get(key string) (string, error)
	Delete(key string) error
}

type Storage struct {
	engine Engine
	logger *zap.Logger
}

func New(e Engine, l *zap.Logger) *Storage {
	return &Storage{
		engine: e,
		logger: l,
	}
}

func (s *Storage) Execute(command *parser.Command) (string, error) {
	err := command.Validate()
	if err != nil {
		return "", fmt.Errorf("failed command.Validate, %w", err)
	}

	s.logger.Info("Executing command", zap.String("action", command.Action), zap.Strings("args", command.Args))
	switch command.Action {
	case GET:
		key := command.Args[0]
		var value string
		value, err = s.engine.Get(key)
		if err != nil {
			return "", fmt.Errorf("failed s.engine.Get, err: %w", err)
		}

		return value, nil
	case SET:
		s.engine.Set(command.Args[0], command.Args[1])
		return "", nil
	case DELETE:
		key := command.Args[0]
		if err = s.engine.Delete(key); err != nil {
			return "", fmt.Errorf("failed s.engine.Delete, err: %w", err)
		}

		return "", nil
	default:
		return "", fmt.Errorf("unknown command: %s", command.Action)
	}
}
