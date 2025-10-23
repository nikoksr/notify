package homeassistant

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	service := New("homeassistant.local:8123", "test-token")
	assert.NotNil(t, service)
	assert.Equal(t, "https://homeassistant.local:8123", service.serverURL)
	assert.Equal(t, "test-token", service.accessToken)
	assert.NotNil(t, service.client)
}

func TestNewWithHTTP(t *testing.T) {
	service := New("http://192.168.1.100:8123", "test-token")
	assert.Equal(t, "http://192.168.1.100:8123", service.serverURL)
}

func TestSendNoToken(t *testing.T) {
	service := New("homeassistant.local:8123", "")

	err := service.Send(context.Background(), "Test", "Message")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "access token is required")
}

func TestSendSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer test-token")
		assert.Contains(t, r.URL.Path, "/api/services/notify/persistent_notification")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := New(server.URL, "test-token")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	assert.NoError(t, err)
}

func TestSendFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	service := New(server.URL, "test-token")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "home assistant API returned status 401")
}

func TestSetTarget(_ *testing.T) {
	service := New("homeassistant.local:8123", "test-token")

	// Should not panic
	service.SetTarget("mobile_app_phone")
	service.SetTarget("notify.email")
}
