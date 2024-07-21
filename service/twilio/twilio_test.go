package twilio

import (
	"context"
	"errors"
	"testing"

	"github.com/kevinburke/twilio-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		fromPhoneNumber string
		toPhoneNumbers  []string
		subject         string
		message         string
		mockSetup       func(*mocktwilioClient)
		expectedError   string
	}{
		{
			name:            "Successful send to single recipient",
			fromPhoneNumber: "+1234567890",
			toPhoneNumbers:  []string{"+0987654321"},
			subject:         "Test Subject",
			message:         "Test Message",
			mockSetup: func(m *mocktwilioClient) {
				m.On("SendMessage", "+1234567890", "+0987654321", "Test Subject\nTest Message", mock.Anything).
					Return(&twilio.Message{}, nil)
			},
			expectedError: "",
		},
		{
			name:            "Successful send to multiple recipients",
			fromPhoneNumber: "+1234567890",
			toPhoneNumbers:  []string{"+0987654321", "+1122334455"},
			subject:         "Test Subject",
			message:         "Test Message",
			mockSetup: func(m *mocktwilioClient) {
				m.On("SendMessage", "+1234567890", mock.AnythingOfType("string"), "Test Subject\nTest Message", mock.Anything).
					Return(&twilio.Message{}, nil).
					Twice()
			},
			expectedError: "",
		},
		{
			name:            "Twilio client error",
			fromPhoneNumber: "+1234567890",
			toPhoneNumbers:  []string{"+0987654321"},
			subject:         "Test Subject",
			message:         "Test Message",
			mockSetup: func(m *mocktwilioClient) {
				m.On("SendMessage", "+1234567890", "+0987654321", "Test Subject\nTest Message", mock.Anything).
					Return(nil, errors.New("Twilio error"))
			},
			expectedError: "send message to recipient \"+0987654321\": Twilio error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mocktwilioClient)
			tt.mockSetup(mockClient)

			s := &Service{
				client:          mockClient,
				fromPhoneNumber: tt.fromPhoneNumber,
				toPhoneNumbers:  tt.toPhoneNumbers,
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
