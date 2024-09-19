// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	parser "github.com/patyukin/mdb/internal/compute/parser"
	mock "github.com/stretchr/testify/mock"
)

// Parser is an autogenerated mock type for the Parser type
type Parser struct {
	mock.Mock
}

// Parse provides a mock function with given fields: _a0
func (_m *Parser) Parse(_a0 string) (*parser.Command, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Parse")
	}

	var r0 *parser.Command
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*parser.Command, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(string) *parser.Command); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parser.Command)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewParser creates a new instance of Parser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewParser(t interface {
	mock.TestingT
	Cleanup(func())
}) *Parser {
	mock := &Parser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
