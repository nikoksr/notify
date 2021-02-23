package whatsapp

import (
	"context"
	"testing"

	"github.com/Rhymen/go-whatsapp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
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

func TestSend(t *testing.T) {
	assert := require.New(t)

	svc := &Service{
		contacts: []string{},
	}

	// test whatsapp client returning error
	mockClient := new(mockWhatsappClient)
	mockClient.On("Send", whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "Contact1@s.whatsapp.net",
		},
		Text: "subject\nmessage",
	}).Return("", errors.New("some error"))
	svc.client = mockClient
	svc.AddReceivers("Contact1")
	ctx := context.Background()
	err := svc.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// test success and multiple receivers
	mockClient = new(mockWhatsappClient)
	mockClient.On("Send", whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "Contact1@s.whatsapp.net",
		},
		Text: "subject\nmessage",
	}).Return("", nil)
	mockClient.On("Send", whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: "Contact2@s.whatsapp.net",
		},
		Text: "subject\nmessage",
	}).Return("", nil)
	svc.client = mockClient
	svc.AddReceivers("Contact1", "Contact2")
	err = svc.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}
