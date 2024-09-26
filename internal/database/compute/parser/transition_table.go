package parser

import (
	"fmt"
	"unicode"
)

// State - текущее состояние
type State int

const (
	StateStart State = iota
	StateCommand
	StateArguments
	StateEnd
)

// Transition условие перехода между состояниями
type Transition struct {
	Condition func(rune) bool
	NextState State
	Action    func(*FSM, rune) error
}

// TransitionTable - таблица переходов для каждого State
var TransitionTable = map[State][]Transition{
	StateStart: {
		{
			Condition: unicode.IsSpace,
			NextState: StateStart,
			Action: func(fsm *FSM, ch rune) error {
				fsm.position++
				return nil
			},
		},
		{
			Condition: isUppercase,
			NextState: StateCommand,
			Action: func(fsm *FSM, ch rune) error {
				fsm.currentToken.WriteRune(ch)
				fsm.position++
				return nil
			},
		},
		{
			Condition: func(r rune) bool {
				return true
			},
			NextState: StateEnd,
			Action: func(fsm *FSM, ch rune) error {
				return fmt.Errorf("invalid character at start: '%c'", ch)
			},
		},
	},
	StateCommand: {
		{
			Condition: unicode.IsSpace,
			NextState: StateArguments,
			Action: func(fsm *FSM, ch rune) error {
				if fsm.currentToken.Len() == 0 {
					return fmt.Errorf("empty command")
				}

				fsm.tokens = append(fsm.tokens, fsm.currentToken.String())
				fsm.currentToken.Reset()
				fsm.position++
				return nil
			},
		},
		{
			Condition: isUppercase,
			NextState: StateCommand,
			Action: func(fsm *FSM, ch rune) error {
				fsm.currentToken.WriteRune(ch)
				fsm.position++
				return nil
			},
		},
		{
			Condition: func(r rune) bool {
				return true
			},
			NextState: StateEnd,
			Action: func(fsm *FSM, ch rune) error {
				return fmt.Errorf("invalid character in command: '%c'", ch)
			},
		},
	},
	StateArguments: {
		{
			Condition: unicode.IsSpace,
			NextState: StateArguments,
			Action: func(fsm *FSM, ch rune) error {
				if fsm.currentToken.Len() > 0 {
					fsm.tokens = append(fsm.tokens, fsm.currentToken.String())
					fsm.currentToken.Reset()
				}

				fsm.position++
				return nil
			},
		},
		{
			Condition: isValidArgumentChar,
			NextState: StateArguments,
			Action: func(fsm *FSM, ch rune) error {
				fsm.currentToken.WriteRune(ch)
				fsm.position++
				if len(fsm.tokens) >= 3 {
					return fmt.Errorf("too many arguments")
				}
				return nil
			},
		},
		{
			Condition: func(r rune) bool {
				return true
			},
			NextState: StateEnd,
			Action: func(fsm *FSM, ch rune) error {
				return fmt.Errorf("invalid character in argument: '%c'", ch)
			},
		},
	},
}
