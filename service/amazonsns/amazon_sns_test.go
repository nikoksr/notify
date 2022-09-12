package amazonsns

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAmazonSNS_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("", "", "")
	assert.NotNil(service)
	assert.Nil(err)
}

func TestAmazonSNS_AddReceivers(t *testing.T) {
	t.Parallel()

	amazonSNS, err := New("", "", "")
	if err != nil {
		t.Error(err)
	}
	amazonSNS.AddReceivers("One topic")
}

func TestAmazonSNS_SendMessageWithNoTopicsConfigured(t *testing.T) {
	t.Parallel()

	mockSns := new(mockSnsSendMessageAPI)
	amazonSNS := AmazonSNS{
		sendMessageClient: mockSns,
	}

	err := amazonSNS.Send(context.Background(), "Subject", "Message")
	assert.Nil(t, err)
	mockSns.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestAmazonSNS_SendMessageWithSucessAndOneTopicConfigured(t *testing.T) {
	t.Parallel()

	mockSns := new(mockSnsSendMessageAPI)
	output := sns.PublishOutput{}
	mockSns.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, nil)

	amazonSNS := AmazonSNS{
		sendMessageClient: mockSns,
	}
	amazonSNS.AddReceivers("arn:aws:sns:region:number:topicname")
	err := amazonSNS.Send(context.Background(), "Subject", "Message")
	assert.Nil(t, err)

	mockSns.AssertExpectations(t)
	assert.Equal(t, 1, len(mockSns.Calls))
}

func TestAmazonSNS_SendMessageWithSucessAndTwoTopicsConfigured(t *testing.T) {
	t.Parallel()

	mockSns := new(mockSnsSendMessageAPI)
	output := sns.PublishOutput{}
	mockSns.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, nil)

	amazonSNS := AmazonSNS{
		sendMessageClient: mockSns,
	}

	amazonSNS.AddReceivers("arn:aws:sns:region:number:topicname1",
		"arn:aws:sns:region:number:topicname1")

	err := amazonSNS.Send(context.Background(), "Subject", "Message")
	assert.Nil(t, err)

	mockSns.AssertExpectations(t)
	assert.Equal(t, 2, len(mockSns.Calls))
}

func TestAmazonSNS_SendMessageWithErrorAndOneQueueConfiguredShouldReturnError(t *testing.T) {
	t.Parallel()

	mockSns := new(mockSnsSendMessageAPI)
	output := sns.PublishOutput{}
	mockSns.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, errors.New("Error on SNS"))

	amazonSNS := AmazonSNS{
		sendMessageClient: mockSns,
	}

	amazonSNS.AddReceivers("arn:aws:sns:region:number:topicname")
	err := amazonSNS.Send(context.Background(), "Subject", "Message")

	assert.NotNil(t, err)

	mockSns.AssertExpectations(t)
	assert.Equal(t, 1, len(mockSns.Calls))
}

func TestAmazonSNS_SendMessageWithErrorAndTwoQueueConfiguredShouldReturnErrorOnFirst(t *testing.T) {
	t.Parallel()

	mockSns := new(mockSnsSendMessageAPI)
	output := sns.PublishOutput{}
	mockSns.On("SendMessage", mock.Anything, mock.Anything, mock.Anything).
		Return(&output, errors.New("Error on SNS"))

	amazonSNS := AmazonSNS{
		sendMessageClient: mockSns,
	}

	amazonSNS.AddReceivers(
		"arn:aws:sns:region:number:topicname1",
		"arn:aws:sns:region:number:topicname1")

	err := amazonSNS.Send(context.Background(), "Subject", "Message")
	assert.NotNil(t, err)
	mockSns.AssertExpectations(t)
	assert.Equal(t, 1, len(mockSns.Calls))
}
