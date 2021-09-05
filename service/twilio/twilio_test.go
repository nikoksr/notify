package twilio

import (
	"context"
	"errors"
	"github.com/kevinburke/twilio-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddReceivers(t *testing.T) {
	assert := require.New(t)

	svc := &Service{
		contacts: []string{},
	}
	contacts := []string{"Contact1", "Contact2", "Contact3"}
	svc.AddReceivers(contacts...)

	assert.Equal(svc.contacts, contacts)
}

func TestSendMessage(t *testing.T) {
	assert := require.New(t)

	svc := &Service{
		contacts: []string{},
		phone:    "From",
	}

	// test twilio client returning error
	mockClient := new(mockTwilioClient)
	call := mockClient.On("SendMessage", "From", "Contact1", "subject\nmessageClient", nil)
	call.Return(&twilio.Message{}, errors.New("some error"))
	svc.messageClient = mockClient
	svc.AddReceivers("Contact1")
	ctx := context.Background()
	err := svc.Send(ctx, "subject", "messageClient")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// test success and multiple receivers
	mockClient = new(mockTwilioClient)
	mockClient.On("SendMessage", "From", "Contact1", "subject\nmessageClient", nil).
		Return(&twilio.Message{}, nil)
	mockClient.On("SendMessage", "From", "Contact2", "subject\nmessageClient", nil).
		Return(&twilio.Message{}, nil)
	svc.messageClient = mockClient
	svc.AddReceivers("Contact1", "Contact2")
	err = svc.Send(ctx, "subject", "messageClient")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
