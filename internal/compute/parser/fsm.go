package parser

import (
	"fmt"
	"github.com/patyukin/mdb/pkg/utils"
	"unicode/utf8"
)

const eof = -1

type fsm struct {
	input string
	pos   int
	start int
	width int
	items []string
	err   error
}

type StateFn func() StateFn

func newFSM(input string) *fsm {
	return &fsm{
		input: input,
	}
}

func (f *fsm) Parse() ([]string, error) {
	for state := f.stateStart; state != nil; {
		state = state()
	}

	if f.err != nil {
		return nil, f.err
	}

	if len(f.items) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	return f.items, nil
}

func (f *fsm) stateStart() StateFn {
	f.skipWhitespace()
	if f.pos >= len(f.input) {
		f.err = fmt.Errorf("empty command")
		return nil
	}

	return f.stateAction
}

func (f *fsm) stateAction() StateFn {
	action := f.consumeToken()
	if f.err != nil {
		return nil
	}

	f.items = append(f.items, action)
	return f.stateArgs
}

func (f *fsm) stateArgs() StateFn {
	for {
		f.skipWhitespace()
		if f.pos >= len(f.input) {
			return nil
		}

		arg := f.consumeToken()
		if arg == "" && f.err == nil {
			f.err = fmt.Errorf("empty argument")
			return nil
		}

		f.items = append(f.items, arg)
	}
}

func (f *fsm) consumeToken() string {
	f.skipWhitespace()
	f.start = f.pos
	for {
		r := f.next()
		if r == eof {
			break
		}

		if utils.IsWhitespace(r) {
			f.backup()
			break
		}

		if !utils.IsValidArgChar(r) {
			f.err = fmt.Errorf("invalid character: %q", r)
			return ""
		}
	}

	return f.input[f.start:f.pos]
}
func (f *fsm) next() rune {
	if f.pos >= len(f.input) {
		f.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(f.input[f.pos:])
	if r == utf8.RuneError && w == 1 {
		f.err = fmt.Errorf("invalid UTF-8 encoding at position %d", f.pos)
		return eof
	}

	f.width = w
	f.pos += f.width

	return r
}

func (f *fsm) backup() {
	f.pos -= f.width
}

func (f *fsm) skipWhitespace() {
	for {
		r := f.peek()
		if !utils.IsWhitespace(r) {
			break
		}

		f.next()
	}
}

func (f *fsm) peek() rune {
	if f.pos >= len(f.input) {
		return eof
	}

	return rune(f.input[f.pos])
}
