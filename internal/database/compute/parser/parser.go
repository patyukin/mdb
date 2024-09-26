package parser

import (
	"fmt"
)

const (
	GET    = "GET"
	SET    = "SET"
	DELETE = "DEL"
)

//go:generate go run github.com/vektra/mockery/v2@v2.45.1 --name=FMS --output ../../mocks
type FMS interface {
	Tokenize() ([]string, error)
}

type Command struct {
	Action string
	Args   []string
}

func (c *Command) Validate() error {
	if c.Action != GET && c.Action != SET && c.Action != DELETE {
		return fmt.Errorf("unknown command: %s", c.Action)
	}

	if (c.Action == GET || c.Action == DELETE) && len(c.Args) != 1 {
		return fmt.Errorf("command %s requires 1 argument", c.Action)
	}

	if c.Action == SET && len(c.Args) != 2 {
		return fmt.Errorf("2 arguments required for SET command")
	}

	return nil
}

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) (*Command, error) {
	f := NewFSM(input)
	currentFSM, err := f.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("failed f.Parse: %w", err)
	}

	cmd := Command{
		Action: currentFSM[0],
		Args:   currentFSM[1:],
	}

	err = cmd.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed cmd.Validate: %w", err)
	}

	return &cmd, nil
}
