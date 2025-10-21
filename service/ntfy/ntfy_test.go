package ntfy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	service := New("ntfy.sh")
	assert.NotNil(t, service)
	assert.Equal(t, "https://ntfy.sh", service.serverURL)
	assert.NotNil(t, service.client)
	assert.Empty(t, service.topics)
}

func TestNewWithHTTP(t *testing.T) {
	service := New("http://localhost:8080")
	assert.Equal(t, "http://localhost:8080", service.serverURL)
}

func TestAddReceivers(t *testing.T) {
	service := New("ntfy.sh")
	service.AddReceivers("topic1", "topic2")
	
	assert.Len(t, service.topics, 2)
	assert.Contains(t, service.topics, "topic1")
	assert.Contains(t, service.topics, "topic2")
}

func TestSendNoTopics(t *testing.T) {
	service := New("ntfy.sh")
	
	err := service.Send(context.Background(), "Test", "Message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no topics configured")
}

func TestSendSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := New(server.URL)
	service.AddReceivers("test-topic")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	assert.NoError(t, err)
}

func TestSendFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := New(server.URL)
	service.AddReceivers("test-topic")

	err := service.Send(context.Background(), "Test Subject", "Test Message")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ntfy API returned status 500")
}