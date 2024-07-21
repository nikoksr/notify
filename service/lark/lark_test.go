package lark

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomAppService_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		receiveIDs    []*ReceiverID
		subject       string
		message       string
		mockSetup     func(*mocksendToer)
		expectedError string
	}{
		{
			name:       "Successful send to single receiver",
			receiveIDs: []*ReceiverID{OpenID("user1")},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mocksendToer) {
				m.On("SendTo", "Test Subject", "Test Message", "user1", "open_id").
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "Successful send to multiple receivers",
			receiveIDs: []*ReceiverID{OpenID("user1"), UserID("user2")},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mocksendToer) {
				m.On("SendTo", "Test Subject", "Test Message", "user1", "open_id").
					Return(nil)
				m.On("SendTo", "Test Subject", "Test Message", "user2", "user_id").
					Return(nil)
			},
			expectedError: "",
		},
		{
			name:       "Lark client error",
			receiveIDs: []*ReceiverID{OpenID("user1")},
			subject:    "Test Subject",
			message:    "Test Message",
			mockSetup: func(m *mocksendToer) {
				m.On("SendTo", "Test Subject", "Test Message", "user1", "open_id").
					Return(errors.New("Lark error"))
			},
			expectedError: "Lark error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSendToer := new(mocksendToer)
			tt.mockSetup(mockSendToer)

			s := &CustomAppService{
				receiveIDs: tt.receiveIDs,
				cli:        mockSendToer,
			}

			err := s.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockSendToer.AssertExpectations(t)
		})
	}
}
