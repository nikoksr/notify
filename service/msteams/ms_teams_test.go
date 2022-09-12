package msteams

import (
	"context"
	"testing"

	teams "github.com/atc0005/go-teams-notify/v2"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestMSTeams_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)
}

func TestMSTeams_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)

	service.AddReceivers("https://outlook.office.com/webhook/...")
	assert.Equal(1, len(service.webHooks))

	service.AddReceivers("https://outlook.office.com/webhook/...", "https://outlook.office.com/webhook/...")
	assert.Equal(3, len(service.webHooks))

	hooks := []string{"https://outlook.office.com/webhook/...", "https://outlook.office.com/webhook/..."}
	service.webHooks = []string{}
	service.AddReceivers(hooks...)

	assert.Equal(service.webHooks, hooks)
}

func TestMSTeams_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)

	// Test no receivers
	ctx := context.Background()
	err := service.Send(ctx, "subject", "message")
	assert.Nil(err)

	// Test error response
	message := teams.NewMessageCard()
	message.Title = "subject"
	message.Text = "message"

	mockClient := newMockTeamsClient(t)
	mockClient.
		On("SendWithContext", ctx, "1234", message).
		Return(errors.New("some error"))

	service.client = mockClient
	service.AddReceivers("1234")
	err = service.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// Test success response
	mockClient = newMockTeamsClient(t)
	mockClient.
		On("SendWithContext", ctx, "1234", message).
		Return(nil)

	mockClient.
		On("SendWithContext", ctx, "5678", message).
		Return(nil)

	service.client = mockClient
	service.AddReceivers("5678")
	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
