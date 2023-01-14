package lark

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLark_NewWebhookService(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := NewWebhookService("")
	assert.NotNil(service)
}

func TestLark_SendWebhook(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	ctx := context.Background()

	// First, test for when the sender returns an error.
	{
		mockSender := newMockSender(t)
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
		mockSender := newMockSender(t)
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
