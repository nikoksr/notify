package dingding

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	cfg := &Config{
		Token:  "test-token",
		Secret: "test-secret",
	}
	service := New(cfg)
	assert.NotNil(t, service)
	assert.NotNil(t, service.client)
	assert.Equal(t, "test-token", service.config.Token)
	assert.Equal(t, "test-secret", service.config.Secret)
}

func TestSend(t *testing.T) {
	cfg := &Config{
		Token:  "invalid-token",
		Secret: "invalid-secret",
	}
	service := New(cfg)
	
	// This will fail with invalid credentials, which is expected
	err := service.Send(context.Background(), "Test Subject", "Test Message")
	assert.Error(t, err)
}