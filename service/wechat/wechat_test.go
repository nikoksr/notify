package wechat

import (
	"context"
	"errors"
	"testing"

	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		userIDs       []string
		subject       string
		content       string
		mockSetup     func(*mockwechatMessageManager)
		expectedError string
	}{
		{
			name:    "Successful send to single user",
			userIDs: []string{"user1"},
			subject: "Test Subject",
			content: "Test Content",
			mockSetup: func(m *mockwechatMessageManager) {
				m.On("Send", mock.MatchedBy(func(msg *message.CustomerMessage) bool {
					return msg.ToUser == "user1" && msg.Msgtype == "text" &&
						msg.Text.Content == "Test Subject\nTest Content"
				})).Return(nil)
			},
			expectedError: "",
		},
		{
			name:    "Successful send to multiple users",
			userIDs: []string{"user1", "user2"},
			subject: "Test Subject",
			content: "Test Content",
			mockSetup: func(m *mockwechatMessageManager) {
				m.On("Send", mock.MatchedBy(func(msg *message.CustomerMessage) bool {
					return (msg.ToUser == "user1" || msg.ToUser == "user2") && msg.Msgtype == "text" &&
						msg.Text.Content == "Test Subject\nTest Content"
				})).Return(nil).Twice()
			},
			expectedError: "",
		},
		{
			name:    "WeChat client error",
			userIDs: []string{"user1"},
			subject: "Test Subject",
			content: "Test Content",
			mockSetup: func(m *mockwechatMessageManager) {
				m.On("Send", mock.Anything).Return(errors.New("WeChat error"))
			},
			expectedError: "send message to user \"user1\": WeChat error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockMessageManager := new(mockwechatMessageManager)
			tt.mockSetup(mockMessageManager)

			s := &Service{
				messageManager: mockMessageManager,
				userIDs:        tt.userIDs,
			}

			err := s.Send(context.Background(), tt.subject, tt.content)

			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			mockMessageManager.AssertExpectations(t)
		})
	}
}
