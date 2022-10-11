package viber

import (
	"context"
	"errors"
	"testing"

	vb "github.com/mileusna/viber"
	"github.com/stretchr/testify/require"
)

func TestViber_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	viber := New("appkey", "senderName", "senderAvatar")
	assert.NotNil(viber)
}

func TestViber_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	viber := New("appkey", "senderName", "senderAvatar")
	assert.Len(viber.SubscribedUserIDs, 0)

	viber.AddReceivers("first-subscriber")
	assert.Len(viber.SubscribedUserIDs, 1)

	viber.AddReceivers("second-subscriber", "third-subscriber")
	assert.Len(viber.SubscribedUserIDs, 3)
}

func TestViber_SetWebhook(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	viber := New("appkey", "senderName", "senderAvatar")
	assert.NotNil(viber)

	// Test error
	viberMock := newMockViberClient(t)
	webhookURLMock := "https://example-webhook.com"
	viberMock.
		On("SetWebhook", webhookURLMock, []string{}).
		Return(vb.WebhookResp{}, errors.New("set webhook error"))

	viber.Client = viberMock
	err := viber.SetWebhook(webhookURLMock)
	assert.NotNil(err)
	viberMock.AssertExpectations(t)

	// Test success
	viberMock = newMockViberClient(t)
	viberMock.
		On("SetWebhook", webhookURLMock, []string{}).
		Return(vb.WebhookResp{}, nil)

	viber.Client = viberMock
	err = viber.SetWebhook(webhookURLMock)
	assert.Nil(err)
	viberMock.AssertExpectations(t)
}

func TestViber_Send(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	viber := New("appkey", "senderName", "senderAvatar")
	assert.NotNil(viber)

	// No receivers added
	ctx := context.Background()
	err := viber.Send(ctx, "subject", "message")
	assert.Nil(err)

	// Test error response
	viberMock := newMockViberClient(t)
	viberMock.
		On("SendTextMessage", "receiver1", "subject\nmessage").
		Return(uint64(0), errors.New("error send text message"))

	viber.Client = viberMock
	viber.AddReceivers("receiver1")
	err = viber.Send(ctx, "subject", "message")
	assert.NotNil(err)
	viberMock.AssertExpectations(t)

	// Test success response
	viberMock = newMockViberClient(t)
	viberMock.
		On("SendTextMessage", "receiver1", "subject\nmessage").
		Return(uint64(0), nil)
	viberMock.
		On("SendTextMessage", "receiver2", "subject\nmessage").
		Return(uint64(0), nil)

	viber.Client = viberMock
	viber.AddReceivers("receiver2")
	err = viber.Send(ctx, "subject", "message")
	assert.Nil(err)
	viberMock.AssertExpectations(t)
}
