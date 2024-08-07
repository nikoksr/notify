// Code generated by mockery v2.43.2. DO NOT EDIT.

package googlechat

import (
	mock "github.com/stretchr/testify/mock"
	chat "google.golang.org/api/chat/v1"
)

// mockspacesMessageCreator is an autogenerated mock type for the spacesMessageCreator type
type mockspacesMessageCreator struct {
	mock.Mock
}

type mockspacesMessageCreator_Expecter struct {
	mock *mock.Mock
}

func (_m *mockspacesMessageCreator) EXPECT() *mockspacesMessageCreator_Expecter {
	return &mockspacesMessageCreator_Expecter{mock: &_m.Mock}
}

// Create provides a mock function with given fields: _a0, _a1
func (_m *mockspacesMessageCreator) Create(_a0 string, _a1 *chat.Message) callCreator {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 callCreator
	if rf, ok := ret.Get(0).(func(string, *chat.Message) callCreator); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(callCreator)
		}
	}

	return r0
}

// mockspacesMessageCreator_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type mockspacesMessageCreator_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - _a0 string
//   - _a1 *chat.Message
func (_e *mockspacesMessageCreator_Expecter) Create(_a0 interface{}, _a1 interface{}) *mockspacesMessageCreator_Create_Call {
	return &mockspacesMessageCreator_Create_Call{Call: _e.mock.On("Create", _a0, _a1)}
}

func (_c *mockspacesMessageCreator_Create_Call) Run(run func(_a0 string, _a1 *chat.Message)) *mockspacesMessageCreator_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(*chat.Message))
	})
	return _c
}

func (_c *mockspacesMessageCreator_Create_Call) Return(_a0 callCreator) *mockspacesMessageCreator_Create_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *mockspacesMessageCreator_Create_Call) RunAndReturn(run func(string, *chat.Message) callCreator) *mockspacesMessageCreator_Create_Call {
	_c.Call.Return(run)
	return _c
}

// newMockspacesMessageCreator creates a new instance of mockspacesMessageCreator. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockspacesMessageCreator(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockspacesMessageCreator {
	mock := &mockspacesMessageCreator{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
