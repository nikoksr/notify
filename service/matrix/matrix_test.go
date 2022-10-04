package matrix

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	matrix "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func TestMatrix_New(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	service, err := New("fake-user-id", "fake-home-server", "fake-home-server", "fake-access-token")
	assert.Nil(err)
	assert.NotNil(service)
	assert.Equal(id.UserID("fake-user-id"), service.options.userID)
	assert.Equal("fake-home-server", service.options.homeServer)
	assert.Equal("fake-access-token", service.options.accessToken)
}

func TestService_Send(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	// Test response
	mockClient := newMockMatrixClient(t)
	mockClient.
		On("SendMessageEvent", id.RoomID("fake-room-id"), event.EventMessage, &Message{Body: "fake-message", Msgtype: event.MsgText}).Return(&matrix.RespSendEvent{}, nil)
	service, _ := New("fake-user-id", "fake-room-id", "fake-home-server", "fake-access-token")
	service.client = mockClient
	err := service.Send(context.Background(), "", "fake-message")
	assert.Nil(err)

	mockClient.AssertExpectations(t)

	// Test error on Send
	mockClient = newMockMatrixClient(t)
	mockClient.
		On("SendMessageEvent", id.RoomID("fake-room-id"), event.EventMessage, &Message{Body: "fake-message", Msgtype: event.MsgText}).Return(nil, errors.New("some-error"))

	service, _ = New("fake-user-id", "fake-room-id", "fake-home-server", "fake-access-token")
	service.client = mockClient
	err = service.Send(context.Background(), "", "fake-message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)
}
