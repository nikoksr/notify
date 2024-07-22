package msteams

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMSTeams_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		webHooks      []string
		subject       string
		message       string
		mockSetup     func(*mockteamsClient)
		expectedError string
	}{
		{
			name:     "Successful send to single webhook",
			webHooks: []string{"https://webhook1.example.com"},
			subject:  "Test Subject",
			message:  "Test Message",
			mockSetup: func(m *mockteamsClient) {
				m.On("SendWithContext", mock.Anything, "https://webhook1.example.com", mock.AnythingOfType("MessageCard")).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:     "Successful send to multiple webhooks",
			webHooks: []string{"https://webhook1.example.com", "https://webhook2.example.com"},
			subject:  "Test Subject",
			message:  "Test Message",
			mockSetup: func(m *mockteamsClient) {
				m.On("SendWithContext", mock.Anything, "https://webhook1.example.com", mock.AnythingOfType("MessageCard")).
					Return(nil)
				m.On("SendWithContext", mock.Anything, "https://webhook2.example.com", mock.AnythingOfType("MessageCard")).
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:     "Teams client error",
			webHooks: []string{"https://webhook1.example.com"},
			subject:  "Test Subject",
			message:  "Test Message",
			mockSetup: func(m *mockteamsClient) {
				m.On("SendWithContext", mock.Anything, "https://webhook1.example.com", mock.AnythingOfType("MessageCard")).
					Return(errors.New("Teams error"))
			},
			expectedError: "send messag to channel \"https://webhook1.example.com\": Teams error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockteamsClient)
			tt.mockSetup(mockClient)

			m := &MSTeams{
				client:   mockClient,
				webHooks: tt.webHooks,
			}

			err := m.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
