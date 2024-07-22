package amazonsns

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAmazonSNS_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		queueTopics   []string
		subject       string
		message       string
		mockSetup     func(*mocksnsSendMessageAPI)
		expectedError string
	}{
		{
			name:        "Successful send to single topic",
			queueTopics: []string{"arn:aws:sns:us-east-1:123456789012:MyTopic"},
			subject:     "Test Subject",
			message:     "Test Message",
			mockSetup: func(m *mocksnsSendMessageAPI) {
				m.On("SendMessage", mock.Anything, mock.AnythingOfType("*sns.PublishInput")).
					Return(&sns.PublishOutput{}, nil)
			},
			expectedError: "",
		},
		{
			name: "Successful send to multiple topics",
			queueTopics: []string{
				"arn:aws:sns:us-east-1:123456789012:MyTopic1",
				"arn:aws:sns:us-east-1:123456789012:MyTopic2",
			},
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mocksnsSendMessageAPI) {
				m.On("SendMessage", mock.Anything, mock.AnythingOfType("*sns.PublishInput")).
					Return(&sns.PublishOutput{}, nil).Times(2)
			},
			expectedError: "",
		},
		{
			name:        "SNS client error",
			queueTopics: []string{"arn:aws:sns:us-east-1:123456789012:MyTopic"},
			subject:     "Test Subject",
			message:     "Test Message",
			mockSetup: func(m *mocksnsSendMessageAPI) {
				m.On("SendMessage", mock.Anything, mock.AnythingOfType("*sns.PublishInput")).
					Return(nil, errors.New("SNS error"))
			},
			expectedError: "send message using Amazon SNS to ARN TOPIC " +
				"\"arn:aws:sns:us-east-1:123456789012:MyTopic\": SNS error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mocksnsSendMessageAPI)
			tt.mockSetup(mockClient)

			s := &AmazonSNS{
				sendMessageClient: mockClient,
				queueTopics:       tt.queueTopics,
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
