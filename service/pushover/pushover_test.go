package pushover

import (
	"context"
	"errors"
	"testing"

	"github.com/gregdel/pushover"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPushover_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		recipients    []pushover.Recipient
		subject       string
		message       string
		mockSetup     func(*mockpushoverClient)
		expectedError string
	}{
		{
			name:       "Successful send to single recipient",
			recipients: []pushover.Recipient{*pushover.NewRecipient("recipient1")},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockpushoverClient) {
				m.On("SendMessage", mock.AnythingOfType("*pushover.Message"), mock.AnythingOfType("*pushover.Recipient")).
					Return(&pushover.Response{}, nil)
			},
			expectedError: "",
		},
		{
			name: "Successful send to multiple recipients",
			recipients: []pushover.Recipient{
				*pushover.NewRecipient("recipient1"),
				*pushover.NewRecipient("recipient2"),
			},
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mockpushoverClient) {
				m.On("SendMessage", mock.AnythingOfType("*pushover.Message"), mock.AnythingOfType("*pushover.Recipient")).
					Return(&pushover.Response{}, nil).
					Times(2)
			},
			expectedError: "",
		},
		{
			name:       "Pushover client error",
			recipients: []pushover.Recipient{*pushover.NewRecipient("recipient1")},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mockpushoverClient) {
				m.On("SendMessage", mock.AnythingOfType("*pushover.Message"), mock.AnythingOfType("*pushover.Recipient")).
					Return(nil, errors.New("Pushover error"))
			},
			expectedError: "send message to recipient 1: Pushover error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockpushoverClient)
			tt.mockSetup(mockClient)

			p := &Pushover{
				client:     mockClient,
				recipients: tt.recipients,
			}

			err := p.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
