package amazonses

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAmazonSES_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("", "", "", "")
	assert.NotNil(service)
	assert.Nil(err)
}

func TestAmazonSES_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("", "", "", "")
	assert.NotNil(service)
	assert.Nil(err)

	receivers := []string{"1", "2", "3", "4", "5"}
	service.AddReceivers(receivers...)

	assert.Equal(service.receiverAddresses, receivers)

	a6 := "6"
	receivers = append(receivers, a6)
	service.AddReceivers(a6)
	assert.Equal(service.receiverAddresses, receivers)
}

func TestAmazonSES_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	sender := "4"
	service, err := New("1", "2", "3", sender)
	assert.NotNil(service)
	assert.Nil(err)

	// Example payload
	input := ses.SendEmailInput{
		Source: &sender,
		Destination: &types.Destination{
			ToAddresses: []string{},
		},
		Message: &types.Message{
			Body: &types.Body{
				Html: &types.Content{
					Data: aws.String("message"),
				},
			},
			Subject: &types.Content{
				Data: aws.String("subject"),
			},
		},
	}

	// No receivers added
	ctx := context.Background()

	mockClient := newMockSesClient(t)
	mockClient.
		On("SendEmail", ctx, &input).
		Return(nil, nil)
	service.client = mockClient

	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// Test with receivers
	receivers := []string{"1", "2"}
	input.Destination.ToAddresses = receivers
	service.AddReceivers(receivers...)

	mockClient = newMockSesClient(t)
	mockClient.
		On("SendEmail", ctx, &input).
		Return(nil, nil)
	service.client = mockClient

	err = service.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// Test with more receivers and error response
	receivers = []string{"1", "2", "3", "4", "5"}
	input.Destination.ToAddresses = receivers
	service.receiverAddresses = make([]string, 0)
	service.AddReceivers(receivers...)

	mockClient = newMockSesClient(t)
	mockClient.
		On("SendEmail", ctx, &input).
		Return(nil, errors.New("some error"))
	service.client = mockClient

	err = service.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)
}
