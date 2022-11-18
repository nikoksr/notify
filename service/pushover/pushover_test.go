package pushover

import (
	"context"
	"testing"

	"github.com/gregdel/pushover"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestPushover_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("")
	assert.NotNil(service)
}

func TestPushover_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("")
	assert.NotNil(service)

	service.AddReceivers("")
	assert.Len(service.recipients, 1)

	service.AddReceivers("", "")
	assert.Len(service.recipients, 3)
}

func TestPushover_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("")
	assert.NotNil(service)

	// No receivers added
	ctx := context.Background()
	err := service.Send(ctx, "subject", "message")
	assert.Nil(err)

	// Test error response
	mockClient := newMockPushoverClient(t)
	mockClient.
		On("SendMessage", &pushover.Message{Title: "subject", Message: "message"}, pushover.NewRecipient("1234")).
		Return(&pushover.Response{}, errors.New("some error"))

	service.client = mockClient
	service.AddReceivers("1234")
	err = service.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// Test success response
	mockClient = newMockPushoverClient(t)
	mockClient.
		On("SendMessage", &pushover.Message{Title: "subject", Message: "message"}, pushover.NewRecipient("1234")).
		Return(&pushover.Response{}, nil)

	mockClient.
		On("SendMessage", &pushover.Message{Title: "subject", Message: "message"}, pushover.NewRecipient("5678")).
		Return(&pushover.Response{}, nil)

	service.client = mockClient
	service.AddReceivers("5678")
	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
