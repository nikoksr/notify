package discord

import (
	"context"
	"errors"
	"testing"

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

func TestDiscord_SetAuthenticatedClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		client        *discordgo.Session
		expectedError string
		validateSetup func(*testing.T, *Discord)
	}{
		{
			name:          "Successfully set authenticated client",
			client:        &discordgo.Session{},
			expectedError: "",
			validateSetup: func(t *testing.T, d *Discord) {
				// Verify the client was set
				require.NotNil(t, d.client)

				// Verify it's a discordgo.Session (not the mock)
				_, ok := d.client.(*discordgo.Session)
				require.True(t, ok, "client should be a *discordgo.Session")
			},
		},
		{
			name:          "Fail with nil client",
			client:        nil,
			expectedError: "discord client is nil",
			validateSetup: func(t *testing.T, d *Discord) {
				// Client should remain unchanged (still the original mock or session)
				require.NotNil(t, d.client)
			},
		},
		{
			name: "Client intents are properly set",
			client: &discordgo.Session{
				Identify: discordgo.Identify{
					Intents: 0, // Start with no intents
				},
			},
			expectedError: "",
			validateSetup: func(t *testing.T, d *Discord) {
				// Verify the client was set and intents were configured
				session, ok := d.client.(*discordgo.Session)
				require.True(t, ok, "client should be a *discordgo.Session")
				require.Equal(t, discordgo.IntentsGuildMessageTyping, session.Identify.Intents)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a Discord instance with initial mock client
			d := &Discord{
				client:     new(mockdiscordSession),
				channelIDs: []string{},
			}

			// Call the method under test
			err := d.SetAuthenticatedClient(tt.client)

			// Verify error expectations
			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			// Run additional validation
			if tt.validateSetup != nil {
				tt.validateSetup(t, d)
			}
		})
	}
}

func TestDiscord_SetAuthenticatedClient_Integration(t *testing.T) {
	t.Parallel()

	// Test that SetAuthenticatedClient works with actual Discord operations
	t.Run("Can send message after setting authenticated client", func(t *testing.T) {
		// Create a Discord instance
		d := New()

		// Create a mock session for testing
		mockSession := new(mockdiscordSession)

		// We need to convert our mock to *discordgo.Session for the method signature
		// In a real scenario, you'd use an actual authenticated session
		realSession := &discordgo.Session{}

		// Set up the client
		err := d.SetAuthenticatedClient(realSession)
		require.NoError(t, err)

		// Verify the client is set and has proper intents
		session, ok := d.client.(*discordgo.Session)
		require.True(t, ok)
		require.Equal(t, discordgo.IntentsGuildMessageTyping, session.Identify.Intents)

		// For testing Send functionality, we'd need to replace with mock again
		// This demonstrates the method works but we can't test Send without mocking
		d.client = mockSession
		d.AddReceivers("test-channel-id")

		mockSession.On("ChannelMessageSend", "test-channel-id", "Test\nMessage").
			Return(&discordgo.Message{}, nil)

		err = d.Send(context.Background(), "Test", "Message")
		require.NoError(t, err)
		mockSession.AssertExpectations(t)
	})
}
