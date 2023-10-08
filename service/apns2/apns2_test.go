package apns2

import (
	"context"
	"testing"

	apnsSvc "github.com/sideshow/apns2"
	"github.com/stretchr/testify/require"
)

func TestAPNS2_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	svc := &Service{
		deviceTokens: []string{},
	}
	deviceTokens := []string{"Token1", "Token2", "Token3"}
	svc.AddReceivers(deviceTokens...)

	assert.Equal(svc.deviceTokens, deviceTokens)
}

func TestBuildNotification(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	notification := buildNotification("<token>", "<topic>", "test notification")
	assert.IsType(new(apnsSvc.Notification), notification)

	assert.Equal(notification.Topic, "<topic>")
	assert.Equal(notification.DeviceToken, "<token>")
}

func TestSend(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	client := &mockApns2Client{}
	notification := buildNotification("<token>", "<topic>", "subject message")
	client.On("Push", notification).Return(&apnsSvc.Response{
		StatusCode: 200,
		Reason:     "",
		ApnsID:     "<apns-id>",
		Timestamp:  apnsSvc.Time{},
	}, nil)

	svc := &Service{
		client:       client,
		topic:        "<topic>",
		deviceTokens: []string{"<token>"},
	}
	err := svc.Send(context.Background(), "subject", "message")

	assert.Nil(err)
}
