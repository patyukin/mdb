package database

import (
	"fmt"
	"github.com/patyukin/mdb/internal/database/compute/parser"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=Storage --output ./mocks
type Storage interface {
	Execute(*parser.Command) (string, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=Compute --output ./mocks
type Compute interface {
	ProcessRequest(request string) (*parser.Command, error)
}

type Database struct {
	strg   Storage
	cmpt   Compute
	logger *zap.Logger
}

func New(cmpt Compute, strg Storage, logger *zap.Logger) *Database {
	return &Database{
		logger: logger,
		strg:   strg,
		cmpt:   cmpt,
	}
}

func (d *Database) HandleQuery(request string) (string, error) {
	cmd, err := d.cmpt.ProcessRequest(request)
	if err != nil {
		return "", fmt.Errorf("failed d.cmpt.ProcessRequest: %w", err)
	}

	result, err := d.strg.Execute(cmd)
	if err != nil {
		return "", fmt.Errorf("failed c.storage.Execute: %w", err)
	}

	d.logger.Info("Request processed successfully")

	return result, nil
}
