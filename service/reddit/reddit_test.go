package reddit

import (
	context "context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func TestReddit_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("id", "secret", "user", "password")
	assert.NotNil(service)
	assert.NoError(err)
}

func TestReddit_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("id", "secret", "user", "password")
	assert.NotNil(service)
	assert.NoError(err)

	service.AddReceivers("")
	assert.Len(service.recipients, 1)

	service.AddReceivers("", "")
	assert.Len(service.recipients, 3)
}

func TestReddit_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("id", "secret", "user", "password")
	assert.NotNil(service)
	assert.NoError(err)

	// No receivers added
	ctx := context.Background()
	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)

	// Test error response
	mockClient := newMockRedditMessageClient(t)
	mockClient.
		On("Send", ctx, mock.AnythingOfType("*reddit.SendMessageRequest")).
		Return(&reddit.Response{}, errors.New("some error"))

	service.client = mockClient
	service.AddReceivers("1234")
	err = service.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// Test success response
	mockClient = newMockRedditMessageClient(t)
	mockClient.
		On("Send", ctx, mock.AnythingOfType("*reddit.SendMessageRequest")).
		Return(&reddit.Response{}, nil)

	service.client = mockClient
	service.AddReceivers("5678")
	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
