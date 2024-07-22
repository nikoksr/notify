package mattermost

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		channelIDs    map[string]bool
		subject       string
		message       string
		mockSetup     func(*mockhttpClient)
		expectedError string
	}{
		{
			name:       "Successful send to single channel",
			channelIDs: map[string]bool{"channel1": true},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockhttpClient) {
				m.On("Send", mock.Anything, "channel1", "Test Subject\nTest Message").
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "Successful send to multiple channels",
			channelIDs: map[string]bool{"channel1": true, "channel2": true},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockhttpClient) {
				m.On("Send", mock.Anything, "channel1", "Test Subject\nTest Message").
					Return(nil)
				m.On("Send", mock.Anything, "channel2", "Test Subject\nTest Message").
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "Mattermost client error",
			channelIDs: map[string]bool{"channel1": true},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockhttpClient) {
				m.On("Send", mock.Anything, "channel1", "Test Subject\nTest Message").
					Return(errors.New("Mattermost error"))
			},
			expectedError: "send message to channel \"channel1\": Mattermost error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockhttpClient)
			tt.mockSetup(mockClient)

			s := &Service{
				messageClient: mockClient,
				channelIDs:    tt.channelIDs,
			}

			err := s.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
