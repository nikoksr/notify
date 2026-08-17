package mailtrap

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
		name     string
		opts     []Option
		wantHost string
		wantPath string
		wantName string
	}{
		{
			name:     "defaults",
			opts:     nil,
			wantHost: defaultSendHost,
			wantPath: sendPath,
			wantName: "",
		},
		{
			name:     "with sender name",
			opts:     []Option{WithSenderName("Sender")},
			wantHost: defaultSendHost,
			wantPath: sendPath,
			wantName: "Sender",
		},
		{
			name:     "with bulk host",
			opts:     []Option{WithBulkHost()},
			wantHost: bulkSendHost,
			wantPath: sendPath,
		},
		{
			name:     "with sandbox",
			opts:     []Option{WithSandbox("123456")},
			wantHost: sandboxHost,
			wantPath: sendPath + "/123456",
		},
		{
			name:     "with custom host",
			opts:     []Option{WithHost("https://example.test")},
			wantHost: "https://example.test",
			wantPath: sendPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			m := New("token", "sender@example.com", tt.opts...)

			assert.Equal(t, "token", m.apiKey)
			assert.Equal(t, "sender@example.com", m.senderAddress)
			assert.Equal(t, tt.wantHost, m.host)
			assert.Equal(t, tt.wantPath, m.path)
			assert.Equal(t, tt.wantName, m.senderName)
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

func TestMailtrap_AddReceivers(t *testing.T) {
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

			require.Len(t, m.receivers, len(tt.want))
			for i, want := range tt.want {
				assert.Equal(t, want, m.receivers[i].Email)
			}
		})
	}
}

func TestMailtrap_BodyFormat(t *testing.T) {
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

func TestMailtrap_Send(t *testing.T) {
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
		assertPath  func(t *testing.T, path string)
	}{
		{
			name:       "html body",
			opts:       []Option{WithSenderName("Sender")},
			receivers:  []string{"to@example.com"},
			statusCode: http.StatusOK,
			assertBody: func(t *testing.T, p payload) {
				t.Helper()
				assert.Equal(t, "sender@example.com", p.From.Email)
				assert.Equal(t, "Sender", p.From.Name)
				require.Len(t, p.To, 1)
				assert.Equal(t, "to@example.com", p.To[0].Email)
				assert.Equal(t, "subject", p.Subject)
				assert.Equal(t, "<b>hello</b>", p.HTML)
				assert.Empty(t, p.Text)
			},
			assertPath: func(t *testing.T, path string) {
				t.Helper()
				assert.Equal(t, sendPath, path)
			},
		},
		{
			name:       "plain text body",
			plainText:  true,
			receivers:  []string{"to@example.com"},
			statusCode: http.StatusOK,
			assertBody: func(t *testing.T, p payload) {
				t.Helper()
				assert.Equal(t, "<b>hello</b>", p.Text)
				assert.Empty(t, p.HTML)
				assert.Empty(t, p.From.Name)
			},
		},
		{
			name:       "multiple receivers",
			receivers:  []string{"a@example.com", "b@example.com"},
			statusCode: http.StatusOK,
			assertBody: func(t *testing.T, p payload) {
				t.Helper()
				require.Len(t, p.To, 2)
				assert.Equal(t, "a@example.com", p.To[0].Email)
				assert.Equal(t, "b@example.com", p.To[1].Email)
			},
		},
		{
			name:        "error status returns error",
			receivers:   []string{"to@example.com"},
			statusCode:  http.StatusUnauthorized,
			respBody:    `{"errors":["Unauthorized"]}`,
			wantErr:     true,
			wantErrText: "401",
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
				assert.Equal(t, "token", r.Header.Get("Api-Token"))
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

			opts := append([]Option{WithHost(server.URL)}, tt.opts...)
			m := New("token", "sender@example.com", opts...)
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
			if tt.assertBody != nil {
				tt.assertBody(t, got)
			}
			if tt.assertPath != nil {
				tt.assertPath(t, gotPath)
			}
		})
	}
}
