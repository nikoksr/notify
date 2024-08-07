// Code generated by mockery v2.43.2. DO NOT EDIT.

package matrix

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	mautrix "maunium.net/go/mautrix"
	event "maunium.net/go/mautrix/event"
	id "maunium.net/go/mautrix/id"
)

// mockmatrixClient is an autogenerated mock type for the matrixClient type
type mockmatrixClient struct {
	mock.Mock
}

type mockmatrixClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockmatrixClient) EXPECT() *mockmatrixClient_Expecter {
	return &mockmatrixClient_Expecter{mock: &_m.Mock}
}

// SendMessageEvent provides a mock function with given fields: ctx, roomID, eventType, contentJSON, extra
func (_m *mockmatrixClient) SendMessageEvent(ctx context.Context, roomID id.RoomID, eventType event.Type, contentJSON interface{}, extra ...mautrix.ReqSendEvent) (*mautrix.RespSendEvent, error) {
	_va := make([]interface{}, len(extra))
	for _i := range extra {
		_va[_i] = extra[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, roomID, eventType, contentJSON)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SendMessageEvent")
	}

	var r0 *mautrix.RespSendEvent
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, id.RoomID, event.Type, interface{}, ...mautrix.ReqSendEvent) (*mautrix.RespSendEvent, error)); ok {
		return rf(ctx, roomID, eventType, contentJSON, extra...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, id.RoomID, event.Type, interface{}, ...mautrix.ReqSendEvent) *mautrix.RespSendEvent); ok {
		r0 = rf(ctx, roomID, eventType, contentJSON, extra...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mautrix.RespSendEvent)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, id.RoomID, event.Type, interface{}, ...mautrix.ReqSendEvent) error); ok {
		r1 = rf(ctx, roomID, eventType, contentJSON, extra...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockmatrixClient_SendMessageEvent_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendMessageEvent'
type mockmatrixClient_SendMessageEvent_Call struct {
	*mock.Call
}

// SendMessageEvent is a helper method to define mock.On call
//   - ctx context.Context
//   - roomID id.RoomID
//   - eventType event.Type
//   - contentJSON interface{}
//   - extra ...mautrix.ReqSendEvent
func (_e *mockmatrixClient_Expecter) SendMessageEvent(ctx interface{}, roomID interface{}, eventType interface{}, contentJSON interface{}, extra ...interface{}) *mockmatrixClient_SendMessageEvent_Call {
	return &mockmatrixClient_SendMessageEvent_Call{Call: _e.mock.On("SendMessageEvent",
		append([]interface{}{ctx, roomID, eventType, contentJSON}, extra...)...)}
}

func (_c *mockmatrixClient_SendMessageEvent_Call) Run(run func(ctx context.Context, roomID id.RoomID, eventType event.Type, contentJSON interface{}, extra ...mautrix.ReqSendEvent)) *mockmatrixClient_SendMessageEvent_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]mautrix.ReqSendEvent, len(args)-4)
		for i, a := range args[4:] {
			if a != nil {
				variadicArgs[i] = a.(mautrix.ReqSendEvent)
			}
		}
		run(args[0].(context.Context), args[1].(id.RoomID), args[2].(event.Type), args[3].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *mockmatrixClient_SendMessageEvent_Call) Return(resp *mautrix.RespSendEvent, err error) *mockmatrixClient_SendMessageEvent_Call {
	_c.Call.Return(resp, err)
	return _c
}

func (_c *mockmatrixClient_SendMessageEvent_Call) RunAndReturn(run func(context.Context, id.RoomID, event.Type, interface{}, ...mautrix.ReqSendEvent) (*mautrix.RespSendEvent, error)) *mockmatrixClient_SendMessageEvent_Call {
	_c.Call.Return(run)
	return _c
}

// newMockmatrixClient creates a new instance of mockmatrixClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockmatrixClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockmatrixClient {
	mock := &mockmatrixClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
