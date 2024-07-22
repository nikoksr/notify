package slack

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSlack_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		channelIDs    []string
		subject       string
		message       string
		mockSetup     func(*mockslackClient)
		expectedError string
	}{
		{
			name:       "Successful send to single channel",
			channelIDs: []string{"C1234567890"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockslackClient) {
				m.On("PostMessageContext", mock.Anything, "C1234567890", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", nil)
			},
			expectedError: "",
		},
		{
			name:       "Successful send to multiple channels",
			channelIDs: []string{"C1234567890", "C0987654321"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockslackClient) {
				m.On("PostMessageContext", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("slack.MsgOption")).
					Return("", "", nil).
					Twice()
			},
			expectedError: "",
		},
		{
			name:       "Slack client error",
			channelIDs: []string{"C1234567890"},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockslackClient) {
				m.On("PostMessageContext", mock.Anything, "C1234567890", mock.AnythingOfType("slack.MsgOption")).
					Return("", "", errors.New("Slack error"))
			},
			expectedError: "send message to channel \"C1234567890\": Slack error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockslackClient)
			tt.mockSetup(mockClient)

			s := &Slack{
				client:     mockClient,
				channelIDs: tt.channelIDs,
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
