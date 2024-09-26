// Code generated by mockery v2.45.1. DO NOT EDIT.

package mocks

import (
	parser "github.com/patyukin/mdb/internal/database/compute/parser"
	mock "github.com/stretchr/testify/mock"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *Storage) Execute(_a0 *parser.Command) (string, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*parser.Command) (string, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(*parser.Command) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*parser.Command) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
