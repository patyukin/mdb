package storage

import (
	"github.com/patyukin/mdb/internal/compute/parser"
	"github.com/patyukin/mdb/pkg/utils"
	"strings"

	"fmt"
	"go.uber.org/zap"
)

const (
	SET    = "SET"
	GET    = "GET"
	DELETE = "DEL"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=Engine --output ../mocks
type Engine interface {
	Set(key string, value string)
	Get(key string) (string, error)
	Delete(key string)
	GetByPattern(pattern string) (map[string]string, error)
	DelByPattern(pattern string) error
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
	s.logger.Info("Executing command", zap.String("action", command.Action), zap.Strings("args", command.Args))
	switch command.Action {
	case GET:
		key := command.Args[0]
		if utils.ContainsWildcard(key) {
			values, err := s.engine.GetByPattern(key)
			if err != nil {
				return "", fmt.Errorf("failed s.engine.GetByPattern, err: %w", err)
			}

			var result strings.Builder
			for k, v := range values {
				result.WriteString(fmt.Sprintf("%s: %s\n", k, v))
			}

			return result.String(), nil
		} else {
			value, err := s.engine.Get(key)
			if err != nil {
				return "", fmt.Errorf("failed s.engine.Get, err: %w", err)
			}

			return value, nil
		}
	case SET:
		s.engine.Set(command.Args[0], command.Args[1])
		return "", nil
	case DELETE:
		key := command.Args[0]
		if utils.ContainsWildcard(key) {
			err := s.engine.DelByPattern(key)
			if err != nil {
				return "", fmt.Errorf("failed s.engine.DelByPattern, err: %w", err)
			}

			return "", nil
		} else {
			s.engine.Delete(key)
			return "", nil
		}
	default:
		return "", fmt.Errorf("unknown command: %s", command.Action)
	}
}
