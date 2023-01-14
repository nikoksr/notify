package slack

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSlack_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("")
	assert.NotNil(service)
}

func TestSlack_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("")
	assert.NotNil(service)

	service.AddReceivers("")
	assert.Len(service.channelIDs, 1)

	service.AddReceivers("", "")
	assert.Len(service.channelIDs, 3)

	service.channelIDs = []string{}
	receivers := []string{"", ""}
	service.AddReceivers(receivers...)

	diff := cmp.Diff(service.channelIDs, receivers)
	assert.Equal("", diff)
}

func TestSlack_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("")
	assert.NotNil(service)

	// No receivers added
	ctx := context.Background()
	err := service.Send(ctx, "subject", "message")
	assert.Nil(err)

	// Test error response
	mockClient := newMockSlackClient(t)
	mockClient.
		On("PostMessageContext", ctx, "1234", mock.AnythingOfType("MsgOption")).
		Return("", "", errors.New("some error"))

	service.client = mockClient
	service.AddReceivers("1234")
	err = service.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// Test success response
	mockClient = newMockSlackClient(t)
	mockClient.
		On("PostMessageContext", ctx, "1234", mock.AnythingOfType("MsgOption")).
		Return("", "", nil)

	mockClient.
		On("PostMessageContext", ctx, "5678", mock.AnythingOfType("MsgOption")).
		Return("", "", nil)

	service.client = mockClient
	service.AddReceivers("5678")
	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
