package parser

import (
	"fmt"
	"strings"
)

// FSM отвечает за токенизацию входной строки
type FSM struct {
	input        string
	position     int
	currentToken strings.Builder
	tokens       []string
	currentState State
}

func NewFSM(input string) *FSM {
	return &FSM{
		input:        input,
		position:     0,
		tokens:       []string{},
		currentState: StateStart,
	}
}

// Tokenize обрабатывает входную строку и генерирует токены, используя таблицу переходов
func (fsm *FSM) Tokenize() ([]string, error) {
	runes := []rune(fsm.input)

	for fsm.currentState != StateEnd && fsm.position <= len(runes) {
		if fsm.position >= len(runes) {
			switch fsm.currentState {
			case StateStart:
				return nil, fmt.Errorf("empty command")
			case StateCommand:
				if fsm.currentToken.Len() == 0 {
					return nil, fmt.Errorf("empty command")
				}
				fsm.tokens = append(fsm.tokens, fsm.currentToken.String())
			case StateArguments:
				if fsm.currentToken.Len() > 0 {
					fsm.tokens = append(fsm.tokens, fsm.currentToken.String())
				}
			default:
				return nil, fmt.Errorf("unknown state: %v", fsm.currentState)
			}
			break
		}

		currentRune := runes[fsm.position]
		transitions, exists := TransitionTable[fsm.currentState]
		if !exists {
			return nil, fmt.Errorf("unknown state: %v", fsm.currentState)
		}

		matched := false
		for _, transition := range transitions {
			if transition.Condition(currentRune) {
				err := transition.Action(fsm, currentRune)
				if err != nil {
					return nil, fmt.Errorf("failed transition.Action: %w", err)
				}

				fsm.currentState = transition.NextState
				matched = true
				break
			}
		}

		if !matched {
			return nil, fmt.Errorf("no transition from state %v for symbol '%c'", fsm.currentState, currentRune)
		}
	}

	if len(fsm.tokens) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	return fsm.tokens, nil
}
