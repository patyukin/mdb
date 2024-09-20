package parser

import (
	"fmt"
	"github.com/patyukin/mdb/pkg/utils"
	"strings"
	"unicode"
)

// FSM отвечает за токенизацию входной строки
type FSM struct {
	input        string
	position     int
	currentToken strings.Builder
	tokens       []string
}

// StateFunc определяет тип функции состояния
type StateFunc func() (StateFunc, error)

// NewFSM создает новый экземпляр FSM
func NewFSM(input string) *FSM {
	return &FSM{
		input:    input,
		position: 0,
		tokens:   []string{},
	}
}

// Tokenize обрабатывает входную строку и генерирует токены, используя функции состояний
func (fsm *FSM) Tokenize() ([]string, error) {
	var err error
	for state := fsm.stateStart; state != nil; {
		state, err = state()
		if err != nil {
			return nil, fmt.Errorf("failed fsm.stateStart: %w", err)
		}
	}

	if len(fsm.tokens) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	return fsm.tokens, nil
}

// stateStart обрабатывает начальное состояние FSM
func (fsm *FSM) stateStart() (StateFunc, error) {
	fsm.skipWhitespace()
	if fsm.position >= len(fsm.input) {
		return nil, fmt.Errorf("empty input")
	}

	return fsm.stateCommandToken, nil
}

// stateCommandToken обрабатывает состояние разбора первого токена (команды)
func (fsm *FSM) stateCommandToken() (StateFunc, error) {
	fsm.skipWhitespace()
	runes := []rune(fsm.input)

	if fsm.position >= len(runes) {
		return nil, fmt.Errorf("empty input")
	}

	for fsm.position < len(runes) {
		ch := runes[fsm.position]

		if unicode.IsSpace(ch) {
			tokenStr := fsm.currentToken.String()

			if tokenStr == "" {
				return nil, fmt.Errorf("empty command")
			}

			fsm.tokens = append(fsm.tokens, tokenStr)
			fsm.currentToken.Reset()

			return fsm.stateArgumentsToken, nil
		}

		if utils.IsUppercase(ch) {
			fsm.currentToken.WriteRune(ch)
			fsm.position++
			continue
		}

		return nil, fmt.Errorf("invalid character in command: '%c'", ch)
	}

	if fsm.currentToken.Len() > 0 {
		tokenStr := fsm.currentToken.String()

		fsm.tokens = append(fsm.tokens, tokenStr)
		fsm.currentToken.Reset()
	}

	return nil, fmt.Errorf("empty arguments")
}

// stateArgumentsToken обрабатывает состояние разбора последующих токенов (аргументов)
func (fsm *FSM) stateArgumentsToken() (StateFunc, error) {
	fsm.skipWhitespace()
	runes := []rune(fsm.input)

	if fsm.position >= len(runes) {
		return nil, fmt.Errorf("empty arguments")
	}

	for fsm.position < len(runes) {
		ch := runes[fsm.position]

		if len(fsm.tokens) == 3 && fsm.currentToken.Len() > 0 {
			return nil, fmt.Errorf("too many arguments")
		}

		if unicode.IsSpace(ch) {
			if fsm.currentToken.Len() > 0 {
				fsm.tokens = append(fsm.tokens, fsm.currentToken.String())
				fsm.currentToken.Reset()
			}

			fsm.position++
			continue
		}

		if utils.IsValidArgumentChar(ch) {
			fsm.currentToken.WriteRune(ch)
			fsm.position++
			continue
		}

		return nil, fmt.Errorf("invalid character in token: %c", ch)
	}

	if fsm.currentToken.Len() > 0 {
		fsm.tokens = append(fsm.tokens, fsm.currentToken.String())
		fsm.currentToken.Reset()
	}

	return nil, nil
}

func (fsm *FSM) skipWhitespace() {
	for fsm.position < len(fsm.input) && unicode.IsSpace(rune(fsm.input[fsm.position])) {
		fsm.position++
	}
}
