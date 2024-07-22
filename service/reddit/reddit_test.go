package reddit

import (
	"context"
	"errors"
	"testing"

	"github.com/caarlos0/go-reddit/v3/reddit"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestReddit_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		recipients    []string
		subject       string
		message       string
		mockSetup     func(*mockredditMessageClient)
		expectedError string
	}{
		{
			name:       "Successful send to single recipient",
			recipients: []string{"user1"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockredditMessageClient) {
				m.On("Send", mock.Anything, mock.AnythingOfType("*reddit.SendMessageRequest")).
					Return(&reddit.Response{}, nil)
			},
			expectedError: "",
		},
		{
			name:       "Successful send to multiple recipients",
			recipients: []string{"user1", "user2"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockredditMessageClient) {
				m.On("Send", mock.Anything, mock.AnythingOfType("*reddit.SendMessageRequest")).
					Return(&reddit.Response{}, nil).Times(2)
			},
			expectedError: "",
		},
		{
			name:       "Reddit client error",
			recipients: []string{"user1"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockredditMessageClient) {
				m.On("Send", mock.Anything, mock.AnythingOfType("*reddit.SendMessageRequest")).
					Return(nil, errors.New("Reddit error"))
			},
			expectedError: "send message to user \"user1\": Reddit error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockredditMessageClient)
			tt.mockSetup(mockClient)

			r := &Reddit{
				client:     mockClient,
				recipients: tt.recipients,
			}

			err := r.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
