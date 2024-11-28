// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	httpwrap "github.com/nsaltun/userapi/pkg/lib/middlewares/httpwrap"
	mock "github.com/stretchr/testify/mock"
)

// HandlerFunc is an autogenerated mock type for the HandlerFunc type
type HandlerFunc struct {
	mock.Mock
}

// Execute provides a mock function with given fields: _a0
func (_m *HandlerFunc) Execute(_a0 *httpwrap.HttpContext) error {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*httpwrap.HttpContext) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewHandlerFunc creates a new instance of HandlerFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHandlerFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *HandlerFunc {
	mock := &HandlerFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}