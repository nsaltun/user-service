// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	middleware "github.com/nsaltun/userapi/pkg/lib/middleware"
	mock "github.com/stretchr/testify/mock"
)

// CustomHandler is an autogenerated mock type for the CustomHandler type
type CustomHandler struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *CustomHandler) Execute(_a0 *middleware.HttpContext) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*middleware.HttpContext) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewCustomHandler creates a new instance of CustomHandler. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCustomHandler(t interface {
	mock.TestingT
	Cleanup(func())
}) *CustomHandler {
	mock := &CustomHandler{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
