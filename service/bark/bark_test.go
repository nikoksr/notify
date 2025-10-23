package bark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	service := New("test-device-key")
	assert.NotNil(t, service)
	assert.Equal(t, "test-device-key", service.deviceKey)
	assert.NotNil(t, service.client)
}

func TestNewWithServers(t *testing.T) {
	service := NewWithServers("test-key", "https://custom.server.com")
	assert.NotNil(t, service)
	assert.Len(t, service.serverURLs, 1)
	assert.Contains(t, service.serverURLs, "https://custom.server.com/")
}

func TestAddReceivers(t *testing.T) {
	service := New("test-key")
	service.AddReceivers("server1.com", "server2.com")
	assert.Len(t, service.serverURLs, 3) // Default + 2 added
}