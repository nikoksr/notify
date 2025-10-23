package zulip

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	service := New("zulip.example.com", "bot@example.com", "api-key")
	assert.NotNil(t, service)
	assert.Equal(t, "https://zulip.example.com", service.serverURL)
	assert.Equal(t, "bot@example.com", service.email)
	assert.Equal(t, "api-key", service.apiKey)
	assert.NotNil(t, service.client)
}

func TestNewWithHTTP(t *testing.T) {
	service := New("http://localhost:9991", "bot@example.com", "api-key")
	assert.Equal(t, "http://localhost:9991", service.serverURL)
}

func TestSendToStreamNoCredentials(t *testing.T) {
	service := New("zulip.example.com", "", "")

	err := service.SendToStream(context.Background(), "general", "Test", "Message")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "email and API key are required")
}

func TestSendToStreamSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.URL.Path, "/api/v1/messages")

		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "bot@example.com", username)
		assert.Equal(t, "api-key", password)

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := New(server.URL, "bot@example.com", "api-key")

	err := service.SendToStream(context.Background(), "general", "Test Topic", "Test Message")
	assert.NoError(t, err)
}

func TestSendDirectMessageSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := New(server.URL, "bot@example.com", "api-key")

	err := service.SendDirectMessage(context.Background(), "user@example.com", "Private message")
	assert.NoError(t, err)
}

func TestSendFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	service := New(server.URL, "bot@example.com", "api-key")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zulip API returned status 401")
}

func TestSendDefaultStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := New(server.URL, "bot@example.com", "api-key")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	assert.NoError(t, err)
}
