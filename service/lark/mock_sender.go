// Code generated by mockery v2.43.2. DO NOT EDIT.

package lark

import mock "github.com/stretchr/testify/mock"

// mocksender is an autogenerated mock type for the sender type
type mocksender struct {
	mock.Mock
}

type mocksender_Expecter struct {
	mock *mock.Mock
}

func (_m *mocksender) EXPECT() *mocksender_Expecter {
	return &mocksender_Expecter{mock: &_m.Mock}
}

// Send provides a mock function with given fields: subject, message
func (_m *mocksender) Send(subject string, message string) error {
	ret := _m.Called(subject, message)

	if len(ret) == 0 {
		panic("no return value specified for Send")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(subject, message)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mocksender_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type mocksender_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - subject string
//   - message string
func (_e *mocksender_Expecter) Send(subject interface{}, message interface{}) *mocksender_Send_Call {
	return &mocksender_Send_Call{Call: _e.mock.On("Send", subject, message)}
}

func (_c *mocksender_Send_Call) Run(run func(subject string, message string)) *mocksender_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *mocksender_Send_Call) Return(_a0 error) *mocksender_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mocksender_Send_Call) RunAndReturn(run func(string, string) error) *mocksender_Send_Call {
	_c.Call.Return(run)
	return _c
}

// newMocksender creates a new instance of mocksender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMocksender(t interface {
	mock.TestingT
	Cleanup(func())
}) *mocksender {
	mock := &mocksender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
