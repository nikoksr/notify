package discord

import (
	"testing"

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
