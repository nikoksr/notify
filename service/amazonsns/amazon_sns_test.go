package amazonsns

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type SNSSendMessageMock struct {
	mock.Mock
}

func (m *SNSSendMessageMock) SendMessage(ctx context.Context,
	params *sns.PublishInput,
	optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestAddReceivers(t *testing.T) {
	amazonSNS, err := New("", "", "")
	if err != nil {
		t.Error(err)
	}
	amazonSNS.AddReceivers("One topic")
}

func TestSendMessageWithNoTopicsConfigured(t *testing.T) {
	mockSns := new(SNSSendMessageMock)
	amazonSNS := AmazonSNS{
		sendMessageClient: mockSns,
	}

	err := amazonSNS.Send(context.Background(), "Subject", "Message")
	assert.Nil(t, err)
	mockSns.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendMessageWithSucessAndOneTopicCOnfigured(t *testing.T) {
	mockSns := new(SNSSendMessageMock)
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

func TestSendMessageWithSucessAndTwoTopicCOnfigured(t *testing.T) {
	mockSns := new(SNSSendMessageMock)
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

func TestSendMessageWithErrorAndOneQueueConfiguredShouldReturnError(t *testing.T) {
	mockSns := new(SNSSendMessageMock)
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

func TestSendMessageWithErrorAndTwoQueueConfiguredShouldReturnErrorOnFirst(t *testing.T) {
	mockSns := new(SNSSendMessageMock)
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
