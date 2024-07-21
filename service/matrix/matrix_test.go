package matrix

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

func TestMatrix_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		roomID        id.RoomID
		subject       string
		message       string
		mockSetup     func(*mockmatrixClient)
		expectedError string
	}{
		{
			name:    "Successful send",
			roomID:  "!roomID:example.com",
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mockmatrixClient) {
				m.On("SendMessageEvent", mock.Anything, id.RoomID("!roomID:example.com"), event.EventMessage, mock.Anything).
					Return(nil, nil)
			},
			expectedError: "",
		},
		{
			name:    "Matrix client error",
			roomID:  "!roomID:example.com",
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mockmatrixClient) {
				m.On("SendMessageEvent", mock.Anything, id.RoomID("!roomID:example.com"), event.EventMessage, mock.Anything).
					Return(nil, errors.New("Matrix error"))
			},
			expectedError: "failed to send message to the room using Matrix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockmatrixClient)
			tt.mockSetup(mockClient)

			m := &Matrix{
				client: mockClient,
				options: ServiceOptions{
					roomID: tt.roomID,
				},
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
