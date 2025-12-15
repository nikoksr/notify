package discord

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/require"
)

func TestDiscord_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		channelIDs    []string
		subject       string
		message       string
		mockSetup     func(*mockdiscordSession)
		expectedError string
	}{
		{
			name:       "Successful send to single channel",
			channelIDs: []string{"123456789"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockdiscordSession) {
				m.On("ChannelMessageSend", "123456789", "Test Subject\nTest Message").
					Return(&discordgo.Message{}, nil)
			},
			expectedError: "",
		},
		{
			name:       "Successful send to multiple channels",
			channelIDs: []string{"123456789", "987654321"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockdiscordSession) {
				m.On("ChannelMessageSend", "123456789", "Test Subject\nTest Message").
					Return(&discordgo.Message{}, nil)
				m.On("ChannelMessageSend", "987654321", "Test Subject\nTest Message").
					Return(&discordgo.Message{}, nil)
			},
			expectedError: "",
		},
		{
			name:       "Discord client error",
			channelIDs: []string{"123456789"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockdiscordSession) {
				m.On("ChannelMessageSend", "123456789", "Test Subject\nTest Message").
					Return(nil, errors.New("Discord error"))
			},
			expectedError: "send message to Discord channel \"123456789\": Discord error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSession := new(mockdiscordSession)
			tt.mockSetup(mockSession)

			d := &Discord{
				client:     mockSession,
				channelIDs: tt.channelIDs,
			}

			err := d.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockSession.AssertExpectations(t)
		})
	}
}

func TestDefaultSession(t *testing.T) {
	t.Parallel()

	session := DefaultSession()

	// Verify it returns a non-nil session
	require.NotNil(t, session, "DefaultSession should return a non-nil session")

	// Verify critical defaults are set (these come from discordgo.New)
	require.NotNil(t, session.Client, "Client should not be nil")
	require.NotNil(t, session.Ratelimiter, "Ratelimiter should not be nil")
	require.Equal(t, 3, session.MaxRestRetries, "MaxRestRetries should be 3")
	require.True(t, session.ShouldRetryOnRateLimit, "ShouldRetryOnRateLimit should be true")
	require.True(t, session.ShouldReconnectOnError, "ShouldReconnectOnError should be true")
	require.Equal(t, 1, session.ShardCount, "ShardCount should be 1")
	require.True(t, session.Compress, "Compress should be true")
}

func TestDiscord_SetClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		setupSession  func() *discordgo.Session
		setupDiscord  func() *Discord
		validateAfter func(*testing.T, *Discord)
	}{
		{
			name: "Set custom session before authentication",
			setupSession: func() *discordgo.Session {
				session := DefaultSession()
				session.ShardID = 5
				session.MaxRestRetries = 10
				return session
			},
			setupDiscord: New,
			validateAfter: func(t *testing.T, d *Discord) {
				session, ok := d.client.(*discordgo.Session)
				require.True(t, ok, "client should be a *discordgo.Session")
				require.Equal(t, 5, session.ShardID, "ShardID should be preserved")
				require.Equal(t, 10, session.MaxRestRetries, "MaxRestRetries should be preserved")
			},
		},
		{
			name: "Set custom session with custom HTTP client",
			setupSession: func() *discordgo.Session {
				session := DefaultSession()
				session.Client = &http.Client{Timeout: 30 * time.Second}
				return session
			},
			setupDiscord: New,
			validateAfter: func(t *testing.T, d *Discord) {
				session, ok := d.client.(*discordgo.Session)
				require.True(t, ok, "client should be a *discordgo.Session")
				require.NotNil(t, session.Client, "HTTP client should not be nil")
				require.Equal(t, 30*time.Second, session.Client.Timeout, "Custom timeout should be preserved")
			},
		},
		{
			name:         "Replace existing session",
			setupSession: DefaultSession,
			setupDiscord: func() *Discord {
				d := New()
				// Set initial custom client
				initialSession := DefaultSession()
				initialSession.ShardID = 1
				d.SetClient(initialSession)
				return d
			},
			validateAfter: func(t *testing.T, d *Discord) {
				session, ok := d.client.(*discordgo.Session)
				require.True(t, ok, "client should be a *discordgo.Session")
				// ShardID should be default (0) from new session, not 1 from initial
				require.Equal(t, 0, session.ShardID, "ShardID should be from new session")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			session := tt.setupSession()
			d := tt.setupDiscord()

			d.SetClient(session)

			tt.validateAfter(t, d)
		})
	}
}

func TestDiscord_AuthenticatePreservesCustomSession(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		setupSession     func() *discordgo.Session
		authMethod       func(*Discord) error
		expectedToken    string
		validatePreserve func(*testing.T, *discordgo.Session)
	}{
		{
			name: "AuthenticateWithBotToken preserves custom HTTP client",
			setupSession: func() *discordgo.Session {
				session := DefaultSession()
				session.Client = &http.Client{Timeout: 45 * time.Second}
				return session
			},
			authMethod: func(d *Discord) error {
				return d.AuthenticateWithBotToken("test-token")
			},
			expectedToken: "Bot test-token",
			validatePreserve: func(t *testing.T, session *discordgo.Session) {
				require.NotNil(t, session.Client, "HTTP client should not be nil")
				require.Equal(t, 45*time.Second, session.Client.Timeout, "Custom timeout should be preserved")
			},
		},
		{
			name: "AuthenticateWithOAuth2Token preserves custom settings",
			setupSession: func() *discordgo.Session {
				session := DefaultSession()
				session.ShardID = 3
				session.ShardCount = 10
				return session
			},
			authMethod: func(d *Discord) error {
				return d.AuthenticateWithOAuth2Token("oauth-token")
			},
			expectedToken: "Bearer oauth-token",
			validatePreserve: func(t *testing.T, session *discordgo.Session) {
				require.Equal(t, 3, session.ShardID, "ShardID should be preserved")
				require.Equal(t, 10, session.ShardCount, "ShardCount should be preserved")
			},
		},
		{
			name: "Multiple custom settings preserved",
			setupSession: func() *discordgo.Session {
				session := DefaultSession()
				session.Client = &http.Client{Timeout: 60 * time.Second}
				session.ShardID = 2
				session.MaxRestRetries = 5
				return session
			},
			authMethod: func(d *Discord) error {
				return d.AuthenticateWithBotToken("multi-token")
			},
			expectedToken: "Bot multi-token",
			validatePreserve: func(t *testing.T, session *discordgo.Session) {
				require.Equal(t, 60*time.Second, session.Client.Timeout, "Timeout should be preserved")
				require.Equal(t, 2, session.ShardID, "ShardID should be preserved")
				require.Equal(t, 5, session.MaxRestRetries, "MaxRestRetries should be preserved")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			session := tt.setupSession()
			d := New()
			d.SetClient(session)

			err := tt.authMethod(d)
			require.NoError(t, err, "Authentication should succeed")

			// Verify token was set correctly
			sessionAfter, ok := d.client.(*discordgo.Session)
			require.True(t, ok, "client should be a *discordgo.Session")
			require.Equal(t, tt.expectedToken, sessionAfter.Token, "Token should be set correctly")
			require.Equal(t, tt.expectedToken, sessionAfter.Identify.Token, "Identify.Token should be set correctly")

			// Verify intents were set
			require.Equal(
				t,
				discordgo.IntentsGuildMessageTyping,
				sessionAfter.Identify.Intents,
				"Intents should be set",
			)

			// Verify custom settings were preserved
			tt.validatePreserve(t, sessionAfter)
		})
	}
}

func TestDiscord_AuthenticateWithoutSetClient(t *testing.T) {
	t.Parallel()

	// Test that authentication still works without calling SetClient
	d := New()
	err := d.AuthenticateWithBotToken("standard-token")
	require.NoError(t, err, "Authentication should succeed without SetClient")

	session, ok := d.client.(*discordgo.Session)
	require.True(t, ok, "client should be a *discordgo.Session")
	require.Equal(t, "Bot standard-token", session.Token, "Token should be set correctly")
	require.Equal(t, discordgo.IntentsGuildMessageTyping, session.Identify.Intents, "Intents should be set")
}
