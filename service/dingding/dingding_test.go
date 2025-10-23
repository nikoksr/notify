package dingding

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	service := New()
	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
}

func TestAddReceivers(t *testing.T) {
	service := New()
	
	service.AddReceivers("webhook1", "webhook2")
	assert.Len(t, service.webhookURLs, 2)
	assert.Contains(t, service.webhookURLs, "webhook1")
	assert.Contains(t, service.webhookURLs, "webhook2")
}

func TestSendNoReceivers(t *testing.T) {
	service := New()
	
	err := service.Send(context.Background(), "Test", "Message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no webhook URLs")
}

func TestSendWithReceivers(t *testing.T) {
	service := New()
	service.AddReceivers("https://oapi.dingtalk.com/robot/send?access_token=test")
	
	// This will fail with network error, but tests the flow
	err := service.Send(context.Background(), "Test", "Message")
	assert.Error(t, err) // Expected network error
}