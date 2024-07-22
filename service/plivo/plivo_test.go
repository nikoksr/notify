package plivo

import (
	"context"
	"errors"
	"testing"

	"github.com/plivo/plivo-go/v7"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Send(t *testing.T) {
	tests := []struct {
		name          string
		destinations  []string
		subject       string
		message       string
		mockSetup     func(*mockplivoMsgClient)
		expectedError string
	}{
		{
			name:         "Successful send to single destination",
			destinations: []string{"+1234567890"},
			subject:      "Test Subject",
			message:      "Test Message",
			mockSetup: func(m *mockplivoMsgClient) {
				m.On("Create", mock.MatchedBy(func(params plivo.MessageCreateParams) bool {
					return params.Src == "Test Source" && params.Dst == "+1234567890" &&
						params.Text == "Test Subject\nTest Message"
				})).Return(&plivo.MessageCreateResponseBody{}, nil)
			},
			expectedError: "",
		},
		{
			name:         "Successful send to multiple destinations",
			destinations: []string{"+1234567890", "+0987654321"},
			subject:      "Test Subject",
			message:      "Test Message",
			mockSetup: func(m *mockplivoMsgClient) {
				m.On("Create", mock.MatchedBy(func(params plivo.MessageCreateParams) bool {
					return params.Src == "Test Source" && params.Dst == "+1234567890<+0987654321" &&
						params.Text == "Test Subject\nTest Message"
				})).Return(&plivo.MessageCreateResponseBody{}, nil)
			},
			expectedError: "",
		},
		{
			name:         "Plivo client error",
			destinations: []string{"+1234567890"},
			subject:      "Test Subject",
			message:      "Test Message",
			mockSetup: func(m *mockplivoMsgClient) {
				m.On("Create", mock.MatchedBy(func(params plivo.MessageCreateParams) bool {
					return params.Src == "Test Source" && params.Dst == "+1234567890" &&
						params.Text == "Test Subject\nTest Message"
				})).Return(nil, errors.New("Plivo error"))
			},
			expectedError: "send SMS to \"+1234567890\": Plivo error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := newMockplivoMsgClient(t)
			tt.mockSetup(mockClient)

			s := &Service{
				client: mockClient,
				mopts: MessageOptions{
					Source: "Test Source",
				},
			}
			s.AddReceivers(tt.destinations...)

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

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		cOpts         *ClientOptions
		mOpts         *MessageOptions
		expectedError string
	}{
		{
			name:          "Valid options",
			cOpts:         &ClientOptions{AuthID: "test-auth-id", AuthToken: "test-auth-token"},
			mOpts:         &MessageOptions{Source: "test-source"},
			expectedError: "",
		},
		{
			name:          "Nil client options",
			cOpts:         nil,
			mOpts:         &MessageOptions{Source: "test-source"},
			expectedError: "client-options cannot be nil",
		},
		{
			name:          "Nil message options",
			cOpts:         &ClientOptions{AuthID: "test-auth-id", AuthToken: "test-auth-token"},
			mOpts:         nil,
			expectedError: "message-options cannot be nil",
		},
		{
			name:          "Empty source in message options",
			cOpts:         &ClientOptions{AuthID: "test-auth-id", AuthToken: "test-auth-token"},
			mOpts:         &MessageOptions{Source: ""},
			expectedError: "source cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tt.cOpts, tt.mOpts)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
