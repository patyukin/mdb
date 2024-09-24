package compute

import (
	"fmt"
	"github.com/patyukin/mdb/internal/database/compute/parser"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=Parser --output ./mocks
type Parser interface {
	Parse(string) (*parser.Command, error)
}

type Compute struct {
	parser Parser
	logger *zap.Logger
}

func New(p Parser, l *zap.Logger) *Compute {
	return &Compute{
		parser: p,
		logger: l,
	}
}

func (c *Compute) ProcessRequest(request string) (*parser.Command, error) {
	c.logger.Info("Received request", zap.String("request", request))
	command, err := c.parser.Parse(request)
	if err != nil {
		return nil, fmt.Errorf("failed c.parser.Parse: %w", err)
	}

	return command, nil
}
