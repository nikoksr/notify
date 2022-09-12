package plivo

import (
	"context"
	"errors"
	"testing"

	"github.com/plivo/plivo-go/v7"
	"github.com/stretchr/testify/require"
)

func TestPlivo_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	// nil ClientOptions
	svc, err := New(nil, &MessageOptions{})
	assert.NotNil(err)
	assert.Nil(svc)

	// nil MessageOptions
	svc, err = New(&ClientOptions{}, nil)
	assert.NotNil(err)
	assert.Nil(svc)

	// empty source
	svc, err = New(&ClientOptions{}, &MessageOptions{})
	assert.NotNil(err)
	assert.Nil(svc)

	// success
	svc, err = New(&ClientOptions{}, &MessageOptions{Source: "12345"})
	assert.Nil(err)
	assert.NotNil(svc)
}

func TestPlivo_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	svc, err := New(&ClientOptions{}, &MessageOptions{Source: "12345"})
	assert.Nil(err)
	assert.NotNil(svc)

	nums := []string{"1", "2", "3", "4", "5"}
	svc.AddReceivers(nums...)

	assert.Equal(svc.destinations, nums)
}

func TestPlivo_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	svc, err := New(&ClientOptions{}, &MessageOptions{Source: "12345"})
	assert.Nil(err)
	assert.NotNil(svc)

	// no receivers added
	ctx := context.Background()
	err = svc.Send(ctx, "message", "test")
	assert.NotNil(err)

	// test plivo client returning error
	mockClient := new(mockPlivoMsgClient)
	mockClient.On("Create", plivo.MessageCreateParams{Src: "12345", Dst: "67890", Text: "message\ntest"}).
		Return(nil, errors.New("some error"))
	svc.client = mockClient
	svc.AddReceivers("67890")
	err = svc.Send(ctx, "message", "test")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// test success and multiple receivers
	mockClient = new(mockPlivoMsgClient)
	mockClient.On("Create", plivo.MessageCreateParams{Src: "12345", Dst: "67890<09876", Text: "message\ntest"}).
		Return(nil, nil)
	svc.client = mockClient
	svc.AddReceivers("09876")
	err = svc.Send(ctx, "message", "test")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
