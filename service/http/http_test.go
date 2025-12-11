package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Allows us to simulate an error returned from the server on a per-request basis.
const headerTestError = "X-Test-Error"

// testServer encapsulates the httptest.Server and provides methods for testing.
type testServer struct {
	server *httptest.Server
}

// newTestServer creates and returns a new testServer.
func newTestServer() *testServer {
	ts := &testServer{}
	ts.server = httptest.NewServer(ts.handler())
	return ts
}

// handler returns the http.HandlerFunc for the test server.
func (ts *testServer) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Header.Get(headerTestError) == "true":
			w.WriteHeader(http.StatusInternalServerError)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}
}

// Close shuts down the test server.
func (ts *testServer) Close() {
	ts.server.Close()
}

// URL returns the URL of the test server.
func (ts *testServer) URL() string {
	return ts.server.URL
}

// Create a custom serializer that will return an error.
type errorSerializer struct{}

// Marshal is a no-op and always returns an error.
func (errorSerializer) Marshal(_ string, _ any) ([]byte, error) {
	return nil, errors.New("error")
}

func TestNew(t *testing.T) {
	t.Parallel()

	s1 := New()
	assert.NotNil(t, s1, "service should not be nil")

	s2 := New()
	assert.NotNil(t, s2, "service should not be nil")
	assert.Equal(t, s1, s2, "services should be equal")
}

func TestService_WithClient(t *testing.T) {
	t.Parallel()

	service := New()
	assert.NotNil(t, service, "service should not be nil")
	assert.NotNil(t, service.client, "client should not be nil")

	// Create a new client
	client := &http.Client{}
	service.WithClient(client)
	assert.Equal(t, client, service.client, "clients should be equal")

	// Nil client should not change the service client
	service.WithClient(nil)
	assert.Equal(t, client, service.client, "clients should be equal")
}

func TestService_AddReceivers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		service *Service
		urls    []string
	}{
		{
			name:    "test case 1",
			service: New(),
			urls: []string{
				"http://localhost:8080",
				"http://localhost:8081",
				"http://localhost:8082",
			},
		},
		{
			name:    "test case 2",
			service: New(),
			urls:    []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.service.AddReceiversURLs(tt.urls...)
			assert.Equal(t, len(tt.urls), len(tt.service.webhooks), "webhooks should be equal")

			for i, hook := range tt.urls {
				assert.Equal(t, hook, tt.service.webhooks[i].URL, "webhooks should be equal")
			}
		})
	}
}

func TestService_Hooks(t *testing.T) {
	t.Parallel()

	// Create a test server for this specific test
	ts := newTestServer()
	defer ts.Close()

	// Set the local server as the receiver
	service := New()
	service.AddReceiversURLs(ts.URL())

	// Constants for the test
	const (
		testSubject = "test subject"
		testMessage = "test message"
	)

	// Add a very simple pre-send hook. We'll check if the header and body are set correctly.
	service.PreSend(func(req *http.Request) error {
		// At this point, the request should be unmodified as this is the first hook. Unmarshal the bodyRaw and check
		// the
		// subject and message.
		bodyRaw, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("read request bodyRaw: %w", err)
		}

		var body map[string]string
		if err = json.Unmarshal(bodyRaw, &body); err != nil {
			return fmt.Errorf("unmarshal request bodyRaw: %w", err)
		}

		// This implicitly checks the correctness of buildDefaultPayload.
		assert.Equal(t, testSubject, body[defaultSubjectKey], "subject should be equal")
		assert.Equal(t, testMessage, body[defaultMessageKey], "message should be equal")

		// Injecting new headers and bodyRaw
		req.Header.Set("X-Test-1", "test-header")
		req.Header.Set("Content-Type", "text/plain")
		req.Body = io.NopCloser(bytes.NewBufferString("test-body"))

		return nil
	})

	// Adding a second pre-send hook. We'll check if the header and body have been correctly modified by the first hook.
	service.PreSend(func(req *http.Request) error {
		// Check the headers
		assert.Equal(t, "test-header", req.Header.Get("X-Test-1"), "header should be equal")
		assert.Equal(t, "text/plain", req.Header.Get("Content-Type"), "header should be equal")

		// Check the bodyRaw
		bodyRaw, err := io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("read request bodyRaw: %w", err)
		}
		assert.Equal(t, "test-body", string(bodyRaw), "body should be equal")

		// Make sure the body is reset to the original value
		req.Body = io.NopCloser(bytes.NewBuffer(bodyRaw))

		// Also, refresh the Content-Length header. This is required because we've modified the bodyRaw and the test
		// would fail otherwise.
		req.ContentLength = int64(len(bodyRaw))

		// Injecting a new header to confirm that consecutive hooks work as expected
		req.Header.Set("X-Test-2", "test-header-2")
		req.Header.Del("X-Test-1")

		return nil
	})

	// Adding a third pre-send hook. We'll check if the header and body have been correctly modified by the first two
	// hooks.
	service.PreSend(func(req *http.Request) error {
		assert.Equal(t, "test-header-2", req.Header.Get("X-Test-2"), "header should be equal")
		assert.Equal(t, "", req.Header.Get("X-Test-1"), "header should be equal")

		// Modifying the headers one last time to verify that the post-send hook works as expected
		req.Header.Set("X-Test-3", "test-header-3")
		req.Header.Del("X-Test-2")

		return nil
	})

	// Add a very simple post-send hook. We'll inject a custom header and return an error, in case the according http
	// header has been set.
	service.PostSend(func(_ *http.Request, res *http.Response) error {
		res.Header.Set("X-Test-1", "test-header")

		return nil
	})

	// Add a second post-send hook. We'll check if the header has been correctly modified by the first hook.
	service.PostSend(func(_ *http.Request, res *http.Response) error {
		assert.Equal(t, "test-header", res.Header.Get("X-Test-1"), "header should be equal")

		// Injecting a new header to confirm that consecutive hooks work as expected
		res.Header.Set("X-Test-2", "test-header-2")
		res.Header.Del("X-Test-1")

		return nil
	})

	// Add a third post-send hook. We'll check if the header has been correctly modified by the first two hooks.
	service.PostSend(func(_ *http.Request, res *http.Response) error {
		assert.Equal(t, "test-header-2", res.Header.Get("X-Test-2"), "header should be equal")
		assert.Equal(t, "", res.Header.Get("X-Test-1"), "header should be equal")

		return nil
	})

	// Sanity check
	assert.Len(t, service.preSendHooks, 3, "preSendHooks should be equal")
	assert.Len(t, service.postSendHooks, 3, "postSendHooks should be equal")

	// Send a notification
	err := service.Send(context.Background(), testSubject, testMessage)
	require.NoError(t, err, "error should be nil")

	// Now, add a new pre-send hook that sets special header that requests the server to return an error. We'll check if
	// the error is correctly returned.
	service.PreSend(func(req *http.Request) error {
		req.Header.Set(headerTestError, "true")

		return nil
	})

	// Send a notification
	err = service.Send(context.Background(), testSubject, testMessage)
	require.Error(t, err, "error should not be nil")

	// Reset the hooks
	service.preSendHooks = make([]PreSendHookFn, 0)
	service.postSendHooks = make([]PostSendHookFn, 0)

	// Add a pre-send hook that returns an error
	service.PreSend(func(_ *http.Request) error {
		return errors.New("test error")
	})

	// Send a notification
	err = service.Send(context.Background(), testSubject, testMessage)
	require.Error(t, err, "error should not be nil")

	// Reset the hooks again and add a post-send hook that returns an error
	service.preSendHooks = make([]PreSendHookFn, 0)

	service.PostSend(func(_ *http.Request, _ *http.Response) error {
		return errors.New("test error")
	})

	// Send a notification
	err = service.Send(context.Background(), testSubject, testMessage)
	require.Error(t, err, "error should not be nil")
}

func TestService_Send(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a test server for this specific test
	ts := newTestServer()
	defer ts.Close()

	// Create service with local server as receiver
	service := New()
	service.AddReceiversURLs(ts.URL())

	// Sending this notification should work without any issues
	err := service.Send(ctx, "test subject", "test message")
	require.NoError(t, err, "error should be nil")

	// Now, let's reset the receivers and set a custom one, specifically requesting for our test server to return an
	// error. This should result in an error.
	service.webhooks = make([]*Webhook, 0)

	header := http.Header{}
	header.Set(headerTestError, "true")

	service.AddReceivers(&Webhook{
		ContentType:  defaultContentType,
		Header:       header,
		Method:       http.MethodPost,
		URL:          ts.URL(),
		BuildPayload: buildDefaultPayload,
	})

	err = service.Send(ctx, "test subject", "test message")
	require.Error(t, err, "error should not be nil")

	// Reset again, add a functioning receiver again for further tests
	service.webhooks = make([]*Webhook, 0)
	service.AddReceiversURLs(ts.URL())

	// Since we won't reset the receivers list again, add a nil receiver to make sure that the service doesn't crash.
	service.AddReceivers(nil)

	err = service.Send(ctx, "test subject", "test message")
	require.NoError(t, err, "error should be nil")

	// Test setting a custom marshaller that always returns an error
	service.Serializer = errorSerializer{}

	err = service.Send(ctx, "test subject", "test message")
	require.Error(t, err, "error should not be nil")

	// Test context cancellation.
	cancel() // Cancel the context

	err = service.Send(ctx, "test subject", "test message")
	require.Error(t, err, "error should not be nil")
}

func Test_newWebhook(t *testing.T) {
	t.Parallel()

	hook1 := newWebhook("https://example.com")
	assert.NotNil(t, hook1, "hook1 should not be nil")

	hook2 := newWebhook("https://example.com")
	assert.NotNil(t, hook2, "hook2 should not be nil")

	assert.NotEqual(t, hook1, hook2, "hooks should not be equal")
}

func TestWebhook_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		hook *Webhook
		want string
	}{
		{
			name: "empty",
			hook: &Webhook{},
			want: "",
		},
		{
			name: "nil",
			hook: nil,
			want: "",
		},
		{
			name: "test case 1",
			hook: newWebhook("https://example.com"),
			want: "POST https://example.com application/json; charset=utf-8",
		},
		{
			name: "test case 2",
			hook: &Webhook{
				Method:      http.MethodGet, // Doesn't have to make sense, but it's just for testing
				URL:         "https://example.com",
				ContentType: "text/plain",
			},
			want: "GET https://example.com text/plain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equalf(t, tt.want, tt.hook.String(), "String() = %v, want %v", tt.hook.String(), tt.want)
		})
	}
}

func Test_defaultMarshaller_Marshal(t *testing.T) {
	t.Parallel()

	type args struct {
		contentType string
		payload     any
	}
	tests := []struct {
		name    string
		args    args
		wantOut []byte
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "test marshal valid json",
			args: args{
				contentType: "application/json",
				payload:     map[string]any{"test": "test"},
			},
			wantOut: []byte(`{"test":"test"}`),
			wantErr: require.NoError,
		},
		{
			name: "test marshal invalid json",
			args: args{
				contentType: "application/json",
				payload:     map[string]any{"test": make(chan int)},
			},
			wantOut: nil,
			wantErr: require.Error,
		},
		{
			name: "test marshal valid text",
			args: args{
				contentType: "text/plain",
				payload:     "test",
			},
			wantOut: []byte("test"),
			wantErr: require.NoError,
		},
		{
			name: "test marshal invalid text",
			args: args{
				contentType: "text/plain",
				payload:     map[string]any{"test": "test"},
			},
			wantOut: nil,
			wantErr: require.Error,
		},
		{
			name: "test marshal invalid content type",
			args: args{
				contentType: "invalid",
				payload:     map[string]any{"test": "test"},
			},
			wantOut: nil,
			wantErr: require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			serializer := defaultMarshaller{}
			gotOut, err := serializer.Marshal(tt.args.contentType, tt.args.payload)

			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}

			assert.Equalf(t, tt.wantOut, gotOut, "Marshal(%v, %v)", tt.args.contentType, tt.args.payload)
		})
	}
}
