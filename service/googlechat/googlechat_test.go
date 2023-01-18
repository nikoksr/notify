package googlechat

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/chat/v1"
	"google.golang.org/api/option"
)

func TestGoogleChat_New(t *testing.T) {
	t.Parallel()
	withCred := option.WithCredentialsFile("example_credentials.json")
	assert := require.New(t)
	service, err := New(withCred)
	assert.Nil(err)
	assert.NotNil(service)
}

func TestGoogleChat_NewWithContext(t *testing.T) {
	t.Parallel()
	withCred := option.WithCredentialsFile("example_credentials.json")
	assert := require.New(t)
	ctx := context.Background()
	service, err := NewWithContext(ctx, withCred)
	assert.Nil(err)
	assert.NotNil(service)
}

func TestGoogleChat_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := &Service{}

	service.AddReceivers("space_a")
	assert.Len(service.spaces, 1)

	service.AddReceivers("space_b", "space_c")
	assert.Len(service.spaces, 3)

	service.spaces = []string{}
	receivers := []string{"space_a", "space_b"}
	service.AddReceivers(receivers...)

	diff := cmp.Diff(service.spaces, receivers)
	assert.Equal("", diff) // assert that there is no difference
}

func TestGoogleChat_Send(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	assert := require.New(t)

	service := &Service{}
	// No receivers added
	err := service.Send(ctx, "subject", "message")
	assert.Nil(err)

	mockMsgCreator := newMockSpacesMessageCreator(t)

	service.messageCreator = mockMsgCreator
	service.AddReceivers("space_a")

	// Test error response
	failedCall := newMockCreateCall(t)
	failedCall.On("Do").Return(nil, errors.New("something happened"))

	mockMsgCreator.
		On("Create", "spaces/space_a", &chat.Message{Text: "subject\nfailure"}).
		Return(failedCall)

	err = service.Send(ctx, "subject", "failure")
	assert.NotNil(err)
	mockMsgCreator.AssertExpectations(t)

	// Test success response
	successCall := newMockCreateCall(t)
	successCall.On("Do").Return(&chat.Message{Text: "subject\nsuccess"}, nil)

	mockMsgCreator.
		On("Create", "spaces/space_a", &chat.Message{Text: "subject\nsuccess"}).
		Return(successCall)

	err = service.Send(ctx, "subject", "success")
	assert.Nil(err)
	mockMsgCreator.AssertExpectations(t)
}
