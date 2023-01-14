package twilio

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	testing "testing"

	twilio "github.com/kevinburke/twilio-go"
	"github.com/stretchr/testify/require"
)

func TestTwilio_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("", "", "")
	assert.NotNil(service)
	assert.Nil(err)
}

func TestTwilio_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	svc := &Service{
		toPhoneNumbers: []string{},
	}
	toPhoneNumbers := []string{"PhoneNumber1", "PhoneNumber2", "PhoneNumber3"}
	svc.AddReceivers(toPhoneNumbers...)

	assert.Equal(svc.toPhoneNumbers, toPhoneNumbers)
}

func TestTwilio_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	svc := &Service{
		fromPhoneNumber: "my_phone_number",
		toPhoneNumbers:  []string{},
	}

	mockPhoneNumber := "recipient_phone_number"
	mockBody := "subject\nmessage"
	mockError := errors.New("some error")

	// test twilio client send
	mockClient := newMockTwilioClient(t)
	mockClient.On("SendMessage",
		svc.fromPhoneNumber,
		mockPhoneNumber,
		mockBody,
		[]*url.URL{}).Return(&twilio.Message{Body: "a response message"}, nil)
	svc.client = mockClient
	svc.AddReceivers(mockPhoneNumber)
	err := svc.Send(context.Background(), "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// test twilio client send returning error
	mockClient = newMockTwilioClient(t)
	mockClient.On("SendMessage",
		svc.fromPhoneNumber,
		mockPhoneNumber,
		mockBody,
		[]*url.URL{}).Return(nil, mockError)
	svc.client = mockClient
	svc.AddReceivers(mockPhoneNumber)
	err = svc.Send(context.Background(), "subject", "message")
	assert.NotNil(err)
	assert.Equal(
		fmt.Sprintf("failed to send message to phone number '%s' using Twilio: %s", mockPhoneNumber, mockError.Error()),
		err.Error())
	mockClient.AssertExpectations(t)

	// test twilio client send multiple receivers
	anotherMockPhoneNumber := "another_recipient_phone_number"
	mockClient = newMockTwilioClient(t)
	mockClient.On("SendMessage",
		svc.fromPhoneNumber,
		mockPhoneNumber,
		mockBody,
		[]*url.URL{}).Return(&twilio.Message{Body: "a response message"}, nil)
	mockClient.On("SendMessage",
		svc.fromPhoneNumber,
		anotherMockPhoneNumber,
		mockBody,
		[]*url.URL{}).Return(&twilio.Message{Body: "a response message"}, nil)
	svc.client = mockClient
	svc.AddReceivers(mockPhoneNumber)
	svc.AddReceivers(anotherMockPhoneNumber)
	err = svc.Send(context.Background(), "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
