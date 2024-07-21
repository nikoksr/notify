package fcm

import (
	"context"
	"errors"
	"testing"

	"firebase.google.com/go/v4/messaging"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		deviceTokens  []string
		subject       string
		message       string
		mockSetup     func(*mockfcmClient)
		expectedError string
	}{
		{
			name:         "Successful send to single device",
			deviceTokens: []string{"token1"},
			subject:      "Test Subject",
			message:      "Test Message",
			mockSetup: func(m *mockfcmClient) {
				m.On("Send", mock.Anything, mock.AnythingOfType("*messaging.Message")).
					Return(&messaging.BatchResponse{}, nil)
			},
			expectedError: "",
		},
		{
			name:         "Successful send to multiple devices",
			deviceTokens: []string{"token1", "token2"},
			subject:      "Test Subject",
			message:      "Test Message",
			mockSetup: func(m *mockfcmClient) {
				m.On("SendMulticast", mock.Anything, mock.AnythingOfType("*messaging.MulticastMessage")).
					Return(&messaging.BatchResponse{}, nil)
			},
			expectedError: "",
		},
		{
			name:         "FCM client error (single device)",
			deviceTokens: []string{"token1"},
			subject:      "Test Subject",
			message:      "Test Message",
			mockSetup: func(m *mockfcmClient) {
				m.On("Send", mock.Anything, mock.AnythingOfType("*messaging.Message")).
					Return(nil, errors.New("FCM error"))
			},
			expectedError: "send message to FCM device with token \"token1\": FCM error",
		},
		{
			name:         "FCM client error (multiple devices)",
			deviceTokens: []string{"token1", "token2"},
			subject:      "Test Subject",
			message:      "Test Message",
			mockSetup: func(m *mockfcmClient) {
				m.On("SendMulticast", mock.Anything, mock.AnythingOfType("*messaging.MulticastMessage")).
					Return(nil, errors.New("FCM error"))
			},
			expectedError: "send multicast message to FCM devices: FCM error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockfcmClient)
			tt.mockSetup(mockClient)

			s := &Service{
				client:       mockClient,
				deviceTokens: tt.deviceTokens,
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
