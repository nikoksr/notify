package googlechat

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/chat/v1"
)

func TestService_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		spaces        []string
		subject       string
		message       string
		mockSetup     func(*mockspacesMessageCreator, *mockcallCreator)
		expectedError string
	}{
		{
			name:    "Successful send to single space",
			spaces:  []string{"space1"},
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mockspacesMessageCreator, c *mockcallCreator) {
				m.On("Create", "spaces/space1", mock.AnythingOfType("*chat.Message")).
					Return(c)
				c.On("Do").Return(&chat.Message{}, nil)
			},
			expectedError: "",
		},
		{
			name:    "Successful send to multiple spaces",
			spaces:  []string{"space1", "space2"},
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mockspacesMessageCreator, c *mockcallCreator) {
				m.On("Create", "spaces/space1", mock.AnythingOfType("*chat.Message")).
					Return(c)
				m.On("Create", "spaces/space2", mock.AnythingOfType("*chat.Message")).
					Return(c)
				c.On("Do").Return(&chat.Message{}, nil).Times(2)
			},
			expectedError: "",
		},
		{
			name:    "Google Chat client error",
			spaces:  []string{"space1"},
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mockspacesMessageCreator, c *mockcallCreator) {
				m.On("Create", "spaces/space1", mock.AnythingOfType("*chat.Message")).
					Return(c)
				c.On("Do").Return(nil, errors.New("Google Chat error"))
			},
			expectedError: "send message to the google chat space \"space1\": Google Chat error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockSpacesMessageCreator := new(mockspacesMessageCreator)
			mockCallCreator := new(mockcallCreator)
			tt.mockSetup(mockSpacesMessageCreator, mockCallCreator)

			s := &Service{
				messageCreator: mockSpacesMessageCreator,
				spaces:         tt.spaces,
			}

			err := s.Send(context.Background(), tt.subject, tt.message)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockSpacesMessageCreator.AssertExpectations(t)
			mockCallCreator.AssertExpectations(t)
		})
	}
}
