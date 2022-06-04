package lark

import (
	"context"
	"errors"
	"testing"

	"github.com/nikoksr/notify/service/lark/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSendWebhook(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// First, test for when the sender returns an error.
	{
		mockSender := mocks.NewSender(t)
		mockSender.
			On("Send", "subject", "message").
			Return(errors.New(""))

		svc := NewWebhookService("")
		svc.cli = mockSender
		err := svc.Send(ctx, "subject", "message")
		assert.NotNil(err)
		mockSender.AssertExpectations(t)
	}

	// Then test for when the sender does not return an error.
	{
		mockSender := mocks.NewSender(t)
		mockSender.
			On("Send", "subject", "message").
			Return(nil)

		svc := NewWebhookService("")
		svc.cli = mockSender
		err := svc.Send(ctx, "subject", "message")
		assert.Nil(err)
		mockSender.AssertExpectations(t)
	}
}
