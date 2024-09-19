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
	Parse() ([]string, error)
	stateAction() StateFn
	consumeToken() string
	next() rune
	skipWhitespace()
	peek() rune
}

type Command struct {
	Action string
	Args   []string
}

type Parser struct{}

func New() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(input string) (*Command, error) {
	f := newFSM(input)
	currentFSM, err := f.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed f.Parse: %w", err)
	}

	action := currentFSM[0]
	args := currentFSM[1:]

	switch action {
	case SET:
		if len(args) != 2 {
			return nil, fmt.Errorf("2 arguments required for SET command")
		}
	case GET, DELETE:
		if len(args) != 1 {
			return nil, fmt.Errorf("command %s requires 1 argument", action)
		}
	default:
		return nil, fmt.Errorf("unknown command: %s", action)
	}

	return &Command{
		Action: action,
		Args:   args,
	}, nil
}
