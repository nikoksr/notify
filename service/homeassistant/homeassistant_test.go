package homeassistant

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	s1, err := New()
	assert.NotNil(t, s1, "service should not be nil")
	assert.Nil(t, err, "error should not exist")

	s2, err := New()
	assert.NotNil(t, s2, "service should not be nil")
	assert.Equal(t, s1, s2, "services should be equal")
	assert.Nil(t, err, "error should not exist")
}

func TestService_Send(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hookId := "very_long_hook_id_to_test_out"
	subject := "test subject"
	message := "test message"

	sent := false
	svr := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/api/webhook/"+hookId) {
			return
		}

		var m map[string]string
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil || m == nil {
			return
		}

		if s, ok := m["subject"]; !ok || s != subject {
			return
		}

		if s, ok := m["message"]; !ok || s != message {
			return
		}

		sent = true
	}))
	defer svr.Close()

	// Create service with local server as receiver
	service, _ := New()
	service.AddWebhook(svr.URL, hookId, "POST")

	// Sending this notification should work without any issues
	err := service.Send(ctx, subject, message)
	assert.NoError(t, err, "error should be nil")
	assert.True(t, sent, "message was not received")
}
