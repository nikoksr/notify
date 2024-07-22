package amazonses

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAmazonSES_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		receivers     []string
		subject       string
		message       string
		mockSetup     func(*mocksesClient)
		expectedError string
	}{
		{
			name:      "Successful send",
			receivers: []string{"test@example.com"},
			subject:   "Test Subject",
			message:   "Test Message",
			mockSetup: func(m *mocksesClient) {
				m.On("SendEmail", mock.Anything, mock.AnythingOfType("*ses.SendEmailInput")).
					Return(&ses.SendEmailOutput{}, nil)
			},
			expectedError: "",
		},
		{
			name:      "SES client error",
			receivers: []string{"test@example.com"},
			subject:   "Test Subject",
			message:   "Test Message",
			mockSetup: func(m *mocksesClient) {
				m.On("SendEmail", mock.Anything, mock.AnythingOfType("*ses.SendEmailInput")).
					Return(nil, errors.New("SES error"))
			},
			expectedError: "send mail using Amazon SES service: SES error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mocksesClient)
			tt.mockSetup(mockClient)

			s := &AmazonSES{
				client:            mockClient,
				senderAddress:     aws.String("sender@example.com"),
				receiverAddresses: tt.receivers,
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
