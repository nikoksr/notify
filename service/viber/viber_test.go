package viber

import (
	"context"
	"errors"
	"testing"

	"github.com/mileusna/viber"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestViber_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		subscribedUserIDs []string
		subject           string
		message           string
		mockSetup         func(*mockviberClient)
		expectedError     string
	}{
		{
			name:              "Successful send to single user",
			subscribedUserIDs: []string{"user1"},
			subject:           "Test Subject",
			message:           "Test Message",
			mockSetup: func(m *mockviberClient) {
				m.On("SendTextMessage", "user1", "Test Subject\nTest Message").
					Return(uint64(1), nil)
			},
			expectedError: "",
		},
		{
			name:              "Successful send to multiple users",
			subscribedUserIDs: []string{"user1", "user2"},
			subject:           "Test Subject",
			message:           "Test Message",
			mockSetup: func(m *mockviberClient) {
				m.On("SendTextMessage", mock.AnythingOfType("string"), "Test Subject\nTest Message").
					Return(uint64(1), nil).Twice()
			},
			expectedError: "",
		},
		{
			name:              "Viber client error",
			subscribedUserIDs: []string{"user1"},
			subject:           "Test Subject",
			message:           "Test Message",
			mockSetup: func(m *mockviberClient) {
				m.On("SendTextMessage", "user1", "Test Subject\nTest Message").
					Return(uint64(0), errors.New("Viber error"))
			},
			expectedError: "send message to user \"user1\": Viber error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockviberClient)
			tt.mockSetup(mockClient)

			v := &Viber{
				Client:            mockClient,
				SubscribedUserIDs: tt.subscribedUserIDs,
			}

			err := v.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func TestViber_SetWebhook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		webhookURL    string
		mockSetup     func(*mockviberClient)
		expectedError string
	}{
		{
			name:       "Successful webhook set",
			webhookURL: "https://example.com/webhook",
			mockSetup: func(m *mockviberClient) {
				m.On("SetWebhook", "https://example.com/webhook", []string{}).
					Return(viber.WebhookResp{Status: 0, StatusMessage: "ok"}, nil)
			},
			expectedError: "",
		},
		{
			name:       "Viber client error",
			webhookURL: "https://example.com/webhook",
			mockSetup: func(m *mockviberClient) {
				m.On("SetWebhook", "https://example.com/webhook", []string{}).
					Return(viber.WebhookResp{}, errors.New("Viber error"))
			},
			expectedError: "Viber error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockviberClient)
			tt.mockSetup(mockClient)

			v := &Viber{
				Client: mockClient,
			}

			err := v.SetWebhook(tt.webhookURL)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockClient.AssertExpectations(t)
		})
	}
}
