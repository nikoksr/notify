package syslog

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		subject       string
		message       string
		mockSetup     func(*mocksyslogWriter)
		expectedError string
	}{
		{
			name:    "Successful send",
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mocksyslogWriter) {
				m.On("Write", []byte("Test Subject: Test Message")).
					Return(28, nil)
			},
			expectedError: "",
		},
		{
			name:    "Syslog writer error",
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mocksyslogWriter) {
				m.On("Write", []byte("Test Subject: Test Message")).
					Return(0, errors.New("Syslog error"))
			},
			expectedError: "write to syslog: Syslog error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockWriter := new(mocksyslogWriter)
			tt.mockSetup(mockWriter)

			s := &Service{
				writer: mockWriter,
			}

			err := s.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockWriter.AssertExpectations(t)
		})
	}
}
