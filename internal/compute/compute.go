package compute

import (
	"bufio"
	"fmt"
	"github.com/patyukin/mdb/internal/compute/parser"
	"os"
	"strings"

	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=Parser --output ../mocks
type Parser interface {
	Parse(string) (*parser.Command, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=Storage --output ../mocks
type Storage interface {
	Execute(*parser.Command) (string, error)
}

type Compute struct {
	parser  Parser
	storage Storage
	logger  *zap.Logger
}

func New(p Parser, s Storage, l *zap.Logger) *Compute {
	return &Compute{
		parser:  p,
		storage: s,
		logger:  l,
	}
}

func (c *Compute) Start() {
	scanner := bufio.NewScanner(os.Stdin)
	c.logger.Info("Database started. Waiting for commands...")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			c.logger.Error("Error reading from input")
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)
		if input == "" {
			c.logger.Info("Empty command received")
			continue
		}

		result, err := c.ProcessRequest(input)
		if err != nil {
			c.logger.Error("failed c.ProcessRequest", zap.Error(err))
		} else if result != "" {
			c.logger.Info("Request processed successfully", zap.String("result", result))
		}
	}

	if err := scanner.Err(); err != nil {
		c.logger.Error("Error reading from input", zap.Error(err))
	}
}

func (c *Compute) ProcessRequest(request string) (string, error) {
	c.logger.Info("Received request", zap.String("request", request))
	command, err := c.parser.Parse(request)
	if err != nil {
		return "", fmt.Errorf("failed c.parser.Parse: %w", err)
	}

	result, err := c.storage.Execute(command)
	if err != nil {
		return "", fmt.Errorf("failed c.storage.Execute: %w", err)
	}

	c.logger.Info("Request processed successfully")
	return result, nil
}
