package nylas

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	n := New("test-api-key", "test-grant-id", "[email protected]", "Test Sender")

	assert.NotNil(t, n)
	assert.Equal(t, "test-api-key", n.apiKey)
	assert.Equal(t, "test-grant-id", n.grantID)
	assert.Equal(t, "[email protected]", n.senderAddress)
	assert.Equal(t, "Test Sender", n.senderName)
	assert.Equal(t, DefaultBaseURL, n.baseURL)
	assert.False(t, n.usePlainText)
	assert.Empty(t, n.receiverAddresses)
	assert.NotNil(t, n.client)
}

func TestNylas_WithBaseURL(t *testing.T) {
	t.Parallel()

	n := New("test-api-key", "test-grant-id", "[email protected]", "Test Sender")
	n.WithBaseURL("https://api.eu.nylas.com")

	assert.Equal(t, "https://api.eu.nylas.com", n.baseURL)
}

func TestNylas_WithHTTPClient(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{}
	n := New("test-api-key", "test-grant-id", "[email protected]", "Test Sender")
	n.WithHTTPClient(mockClient)

	assert.Equal(t, mockClient, n.client)
}

func TestNylas_AddReceivers(t *testing.T) {
	t.Parallel()

	n := New("test-api-key", "test-grant-id", "[email protected]", "Test Sender")

	n.AddReceivers("[email protected]")
	assert.Len(t, n.receiverAddresses, 1)
	assert.Equal(t, "[email protected]", n.receiverAddresses[0])

	n.AddReceivers("[email protected]", "[email protected]")
	assert.Len(t, n.receiverAddresses, 3)
	assert.Equal(t, "[email protected]", n.receiverAddresses[1])
	assert.Equal(t, "[email protected]", n.receiverAddresses[2])
}

func TestNylas_BodyFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		format       BodyType
		expectedHTML bool
	}{
		{
			name:         "HTML format",
			format:       HTML,
			expectedHTML: true,
		},
		{
			name:         "PlainText format",
			format:       PlainText,
			expectedHTML: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			n := New("test-api-key", "test-grant-id", "[email protected]", "Test Sender")
			n.BodyFormat(tt.format)

			if tt.expectedHTML {
				assert.False(t, n.usePlainText)
			} else {
				assert.True(t, n.usePlainText)
			}
		})
	}
}

func TestNylas_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		subject       string
		message       string
		receivers     []string
		setupMock     func(*mockHTTPClient)
		expectedError string
		expectNoError bool
	}{
		{
			name:      "Successful send to single receiver",
			subject:   "Test Subject",
			message:   "<h1>Test Message</h1>",
			receivers: []string{"[email protected]"},
			setupMock: func(m *mockHTTPClient) {
				m.response = &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(strings.NewReader(`{
						"data": {
							"id": "msg-123",
							"object": "message"
						},
						"request_id": "req-456"
					}`)),
				}
				m.err = nil
			},
			expectNoError: true,
		},
		{
			name:      "Successful send to multiple receivers",
			subject:   "Test Subject",
			message:   "Test Message",
			receivers: []string{"[email protected]", "[email protected]"},
			setupMock: func(m *mockHTTPClient) {
				m.response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"data":{"id":"msg-123"},"request_id":"req-456"}`)),
				}
				m.err = nil
			},
			expectNoError: true,
		},
		{
			name:      "No receivers configured",
			subject:   "Test Subject",
			message:   "Test Message",
			receivers: []string{},
			setupMock: func(_ *mockHTTPClient) {
				// No setup needed as we won't reach the HTTP call
			},
			expectedError: "no receivers configured",
		},
		{
			name:      "HTTP client error",
			subject:   "Test Subject",
			message:   "Test Message",
			receivers: []string{"[email protected]"},
			setupMock: func(m *mockHTTPClient) {
				m.response = nil
				m.err = errors.New("connection timeout")
			},
			expectedError: "send request: connection timeout",
		},
		{
			name:      "API error response with proper error structure",
			subject:   "Test Subject",
			message:   "Test Message",
			receivers: []string{"[email protected]"},
			setupMock: func(m *mockHTTPClient) {
				m.response = &http.Response{
					StatusCode: http.StatusBadRequest,
					Body: io.NopCloser(strings.NewReader(`{
						"error": {
							"type": "invalid_request",
							"message": "Invalid grant ID"
						},
						"request_id": "req-789"
					}`)),
				}
				m.err = nil
			},
			expectedError: "nylas api error: Invalid grant ID (type: invalid_request, request_id: req-789)",
		},
		{
			name:      "API error response without proper error structure",
			subject:   "Test Subject",
			message:   "Test Message",
			receivers: []string{"[email protected]"},
			setupMock: func(m *mockHTTPClient) {
				m.response = &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader(`Internal Server Error`)),
				}
				m.err = nil
			},
			expectedError: "nylas api error (status 500): Internal Server Error",
		},
		{
			name:      "Context cancelled",
			subject:   "Test Subject",
			message:   "Test Message",
			receivers: []string{"[email protected]"},
			setupMock: func(m *mockHTTPClient) {
				m.response = nil
				m.err = context.Canceled
			},
			expectedError: "send request: context canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockClient := &mockHTTPClient{}
			tt.setupMock(mockClient)

			n := New("test-api-key", "test-grant-id", "[email protected]", "Test Sender")
			n.WithHTTPClient(mockClient)
			n.AddReceivers(tt.receivers...)

			err := n.Send(context.Background(), tt.subject, tt.message)

			if tt.expectNoError {
				require.NoError(t, err)
			} else if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}

			// Verify that the request was made with correct headers and URL (if applicable)
			if len(tt.receivers) > 0 && mockClient.lastRequest != nil {
				assert.Equal(t, "Bearer test-api-key", mockClient.lastRequest.Header.Get("Authorization"))
				assert.Equal(t, "application/json", mockClient.lastRequest.Header.Get("Content-Type"))
				assert.Equal(t, "application/json", mockClient.lastRequest.Header.Get("Accept"))
				assert.Contains(t, mockClient.lastRequest.URL.String(), "/v3/grants/test-grant-id/messages/send")
			}
		})
	}
}

func TestNylas_Send_RequestBody(t *testing.T) {
	t.Parallel()

	mockClient := &mockHTTPClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"data":{"id":"msg-123"},"request_id":"req-456"}`)),
		},
	}

	n := New("test-api-key", "test-grant-id", "[email protected]", "Test Sender")
	n.WithHTTPClient(mockClient)
	n.AddReceivers("[email protected]", "[email protected]")

	err := n.Send(context.Background(), "Test Subject", "<h1>Test Message</h1>")
	require.NoError(t, err)

	// Verify the request body contains the expected fields
	require.NotNil(t, mockClient.lastRequest)
	bodyBytes, err := io.ReadAll(mockClient.lastRequest.Body)
	require.NoError(t, err)

	bodyStr := string(bodyBytes)
	assert.Contains(t, bodyStr, `"subject":"Test Subject"`)
	// JSON encoding escapes HTML characters, so check for the escaped version
	assert.Contains(t, bodyStr, `"body":"`)
	assert.Contains(t, bodyStr, `Test Message`)
	assert.Contains(t, bodyStr, `"[email protected]"`)
	assert.Contains(t, bodyStr, `"[email protected]"`)
	assert.Contains(t, bodyStr, `"from":[{"email":"[email protected]","name":"Test Sender"}]`)
}

// mockHTTPClient is a simple mock implementation of httpClient for testing.
type mockHTTPClient struct {
	response    *http.Response
	err         error
	lastRequest *http.Request
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Store the request for verification in tests
	// Clone the request body since it can only be read once
	if req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Create a clone for storage
		clonedReq := req.Clone(req.Context())
		clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		m.lastRequest = clonedReq
	} else {
		m.lastRequest = req
	}

	return m.response, m.err
}
