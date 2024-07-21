// Code generated by mockery v2.43.2. DO NOT EDIT.

package wechat

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// mockverificationCallbackFunc is an autogenerated mock type for the verificationCallbackFunc type
type mockverificationCallbackFunc struct {
	mock.Mock
}

type mockverificationCallbackFunc_Expecter struct {
	mock *mock.Mock
}

func (_m *mockverificationCallbackFunc) EXPECT() *mockverificationCallbackFunc_Expecter {
	return &mockverificationCallbackFunc_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: r, verified
func (_m *mockverificationCallbackFunc) Execute(r *http.Request, verified bool) {
	_m.Called(r, verified)
}

// mockverificationCallbackFunc_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type mockverificationCallbackFunc_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - r *http.Request
//   - verified bool
func (_e *mockverificationCallbackFunc_Expecter) Execute(r interface{}, verified interface{}) *mockverificationCallbackFunc_Execute_Call {
	return &mockverificationCallbackFunc_Execute_Call{Call: _e.mock.On("Execute", r, verified)}
}

func (_c *mockverificationCallbackFunc_Execute_Call) Run(run func(r *http.Request, verified bool)) *mockverificationCallbackFunc_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request), args[1].(bool))
	})
	return _c
}

func (_c *mockverificationCallbackFunc_Execute_Call) Return() *mockverificationCallbackFunc_Execute_Call {
	_c.Call.Return()
	return _c
}

func (_c *mockverificationCallbackFunc_Execute_Call) RunAndReturn(run func(*http.Request, bool)) *mockverificationCallbackFunc_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// newMockverificationCallbackFunc creates a new instance of mockverificationCallbackFunc. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockverificationCallbackFunc(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockverificationCallbackFunc {
	mock := &mockverificationCallbackFunc{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
