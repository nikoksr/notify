package gotify

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	service := New("gotify.example.com", "test-token")
	assert.NotNil(t, service)
	assert.Equal(t, "https://gotify.example.com", service.serverURL)
	assert.Equal(t, "test-token", service.appToken)
	assert.NotNil(t, service.client)
}

func TestNewWithHTTP(t *testing.T) {
	service := New("http://localhost:8080", "test-token")
	assert.Equal(t, "http://localhost:8080", service.serverURL)
}

func TestSendNoToken(t *testing.T) {
	service := New("gotify.example.com", "")

	err := service.Send(context.Background(), "Test", "Message")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "app token is required")
}

func TestSendSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.URL.RawQuery, "token=test-token")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := New(server.URL, "test-token")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	assert.NoError(t, err)
}

func TestSendFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := New(server.URL, "test-token")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gotify API returned status 500")
}

func TestSetPriority(_ *testing.T) {
	service := New("gotify.example.com", "test-token")

	// Should not panic
	service.SetPriority(5)
	service.SetPriority(0)
	service.SetPriority(10)
}
