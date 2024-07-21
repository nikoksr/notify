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
