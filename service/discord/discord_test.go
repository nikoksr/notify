package discord

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDiscord_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	assert.NotNil(New())
}

func TestDiscord_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	service := New()
	assert.NotNil(service)

	channels := []string{"1", "2", "3", "4", "5"}
	service.AddReceivers(channels...)
	assert.Equal(service.channelIDs, channels)
}

func TestDiscord_Authenticate(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	service := New()
	assert.NotNil(service)

	// Note: The following might look confusing, because the validation mechanism is not mocked and never returns an
	// error. The function name may be misleading because it is not actually testing the authentication mechanism. The
	// actual authentication only happens when the service is sends a message.

	// OAuth2
	err := service.AuthenticateWithOAuth2Token("12345")
	assert.Nil(err)

	err = service.AuthenticateWithOAuth2Token("")
	assert.Nil(err)

	// Bot token
	err = service.AuthenticateWithBotToken("12345")
	assert.Nil(err)

	err = service.AuthenticateWithBotToken("")
	assert.Nil(err)
}

func TestDiscord_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)

	// No receivers added
	ctx := context.Background()
	err := service.Send(ctx, "subject", "message")
	assert.Nil(err)

	// Test error response
	mockClient := newMockDiscordSession(t)
	mockClient.
		On("ChannelMessageSend", "1234", "subject\nmessage").
		Return(nil, errors.New("some error"))

	service.client = mockClient
	service.AddReceivers("1234")
	err = service.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// Test success response
	mockClient = newMockDiscordSession(t)
	mockClient.
		On("ChannelMessageSend", "1234", "subject\nmessage").
		Return(nil, nil)

	mockClient.
		On("ChannelMessageSend", "5678", "subject\nmessage").
		Return(nil, nil)

	service.client = mockClient
	service.AddReceivers("5678")
	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
