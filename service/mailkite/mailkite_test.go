package mailkite

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        []Option
		wantBaseURL string
		wantName    string
		wantReplyTo string
	}{
		{
			name:        "defaults",
			opts:        nil,
			wantBaseURL: defaultBaseURL,
		},
		{
			name:        "with sender name",
			opts:        []Option{WithSenderName("Sender")},
			wantBaseURL: defaultBaseURL,
			wantName:    "Sender",
		},
		{
			name:        "with reply-to",
			opts:        []Option{WithReplyTo("support@example.com")},
			wantBaseURL: defaultBaseURL,
			wantReplyTo: "support@example.com",
		},
		{
			name:        "with custom base url",
			opts:        []Option{WithBaseURL("https://example.test")},
			wantBaseURL: "https://example.test",
		},
		{
			name:        "with custom base url trailing slash",
			opts:        []Option{WithBaseURL("https://example.test/")},
			wantBaseURL: "https://example.test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := New("mk_live_token", "sender@example.com", tt.opts...)

			assert.Equal(t, "mk_live_token", m.apiKey)
			assert.Equal(t, "sender@example.com", m.senderAddress)
			assert.Equal(t, tt.wantBaseURL, m.baseURL)
			assert.Equal(t, tt.wantName, m.senderName)
			assert.Equal(t, tt.wantReplyTo, m.replyTo)
			assert.False(t, m.usePlainText)
			assert.Empty(t, m.receivers)
		})
	}
}

func TestWithHTTPClient(t *testing.T) {
	t.Parallel()

	t.Run("custom client is used", func(t *testing.T) {
		t.Parallel()

		client := &http.Client{}
		m := New("token", "sender@example.com", WithHTTPClient(client))

		assert.Same(t, client, m.client)
	})

	t.Run("nil client is ignored", func(t *testing.T) {
		t.Parallel()

		m := New("token", "sender@example.com", WithHTTPClient(nil))

		assert.Same(t, http.DefaultClient, m.client)
	})
}

func TestMailKite_AddReceivers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		calls [][]string
		want  []string
	}{
		{
			name:  "single call single address",
			calls: [][]string{{"a@example.com"}},
			want:  []string{"a@example.com"},
		},
		{
			name:  "single call multiple addresses",
			calls: [][]string{{"a@example.com", "b@example.com"}},
			want:  []string{"a@example.com", "b@example.com"},
		},
		{
			name:  "multiple calls accumulate",
			calls: [][]string{{"a@example.com"}, {"b@example.com", "c@example.com"}},
			want:  []string{"a@example.com", "b@example.com", "c@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := New("token", "sender@example.com")
			for _, call := range tt.calls {
				m.AddReceivers(call...)
			}

			assert.Equal(t, tt.want, m.receivers)
		})
	}
}

func TestMailKite_BodyFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		format        BodyType
		wantPlainText bool
	}{
		{name: "html", format: HTML, wantPlainText: false},
		{name: "plain text", format: PlainText, wantPlainText: true},
		{name: "unknown defaults to html", format: BodyType(42), wantPlainText: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := New("token", "sender@example.com")
			m.BodyFormat(tt.format)

			assert.Equal(t, tt.wantPlainText, m.usePlainText)
		})
	}
}

func TestMailKite_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        []Option
		plainText   bool
		receivers   []string
		statusCode  int
		respBody    string
		wantErr     bool
		wantErrText string
		assertBody  func(t *testing.T, p payload)
	}{
		{
			name:       "html body with sender name",
			opts:       []Option{WithSenderName("Sender")},
			receivers:  []string{"to@example.com"},
			statusCode: http.StatusOK,
			respBody:   `{"id":"msg_123","status":"queued"}`,
			assertBody: func(t *testing.T, p payload) {
				t.Helper()
				assert.Equal(t, "Sender <sender@example.com>", p.From)
				assert.Equal(t, []string{"to@example.com"}, p.To)
				assert.Equal(t, "subject", p.Subject)
				assert.Equal(t, "<b>hello</b>", p.HTML)
				assert.Empty(t, p.Text)
				assert.Empty(t, p.ReplyTo)
			},
		},
		{
			name:       "plain text body",
			plainText:  true,
			receivers:  []string{"to@example.com"},
			statusCode: http.StatusOK,
			assertBody: func(t *testing.T, p payload) {
				t.Helper()
				assert.Equal(t, "sender@example.com", p.From)
				assert.Equal(t, "<b>hello</b>", p.Text)
				assert.Empty(t, p.HTML)
			},
		},
		{
			name:       "multiple receivers and reply-to",
			opts:       []Option{WithReplyTo("support@example.com")},
			receivers:  []string{"a@example.com", "b@example.com"},
			statusCode: http.StatusAccepted,
			assertBody: func(t *testing.T, p payload) {
				t.Helper()
				assert.Equal(t, []string{"a@example.com", "b@example.com"}, p.To)
				assert.Equal(t, "support@example.com", p.ReplyTo)
			},
		},
		{
			name:        "error status returns error",
			receivers:   []string{"to@example.com"},
			statusCode:  http.StatusUnauthorized,
			respBody:    `{"error":"invalid api key"}`,
			wantErr:     true,
			wantErrText: "401",
		},
		{
			name:        "error body is included",
			receivers:   []string{"to@example.com"},
			statusCode:  http.StatusForbidden,
			respBody:    `{"error":"domain not verified"}`,
			wantErr:     true,
			wantErrText: "domain not verified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got payload
			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "Bearer mk_live_token", r.Header.Get("Authorization"))
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				b, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				assert.NoError(t, json.Unmarshal(b, &got))

				w.WriteHeader(tt.statusCode)
				if tt.respBody != "" {
					_, _ = w.Write([]byte(tt.respBody))
				}
			}))
			defer server.Close()

			opts := append([]Option{WithBaseURL(server.URL)}, tt.opts...)
			m := New("mk_live_token", "sender@example.com", opts...)
			if tt.plainText {
				m.BodyFormat(PlainText)
			}
			m.AddReceivers(tt.receivers...)

			err := m.Send(context.Background(), "subject", "<b>hello</b>")

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrText)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, sendPath, gotPath)
			if tt.assertBody != nil {
				tt.assertBody(t, got)
			}
		})
	}
}

func TestMailKite_Send_RequestError(t *testing.T) {
	t.Parallel()

	m := New("token", "sender@example.com", WithBaseURL("http://127.0.0.1:1"))
	m.AddReceivers("to@example.com")

	err := m.Send(context.Background(), "subject", "message")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "send message")
}
