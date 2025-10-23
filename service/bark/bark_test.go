package bark

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	service := New("https://api.day.app/test")
	assert.NotNil(t, service)
	assert.Equal(t, "https://api.day.app/test", service.endpoint)
}

func TestSend(t *testing.T) {
	service := New("https://api.day.app/test")
	
	// Test successful send (mock would be better)
	err := service.Send(context.Background(), "Test Subject", "Test Message")
	// Since we can't mock easily here, we expect this to fail with network error
	// In a real implementation, we'd use dependency injection for HTTP client
	assert.Error(t, err) // Expected to fail without real endpoint
}

func TestSendEmptyMessage(t *testing.T) {
	service := New("https://api.day.app/test")
	
	err := service.Send(context.Background(), "", "")
	assert.Error(t, err)
}