package bark

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetEncryptionKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		key     string
		wantErr string
	}{
		{
			name: "AES-128",
			key:  strings.Repeat("k", aesKeyLen128),
		},
		{
			name: "AES-192",
			key:  strings.Repeat("k", aesKeyLen192),
		},
		{
			name: "AES-256",
			key:  strings.Repeat("k", aesKeyLen256),
		},
		{
			name: "empty key disables encryption",
			key:  "",
		},
		{
			name:    "too short",
			key:     "short",
			wantErr: "bark encryption key must contain exactly 16, 24, or 32 ASCII characters",
		},
		{
			name:    "unsupported length",
			key:     strings.Repeat("k", 17),
			wantErr: "bark encryption key must contain exactly 16, 24, or 32 ASCII characters",
		},
		{
			name:    "non-ascii",
			key:     strings.Repeat("锁", 16),
			wantErr: "bark encryption key must contain only ASCII characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := New("device-key")
			err := svc.SetEncryptionKey(tt.key)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)
				assert.Empty(t, svc.encryptionKey)
				return
			}

			require.NoError(t, err)
			if tt.key == "" {
				assert.Nil(t, svc.encryptionKey)
			} else {
				assert.Equal(t, []byte(tt.key), svc.encryptionKey)
			}
		})
	}
}

func TestSetEncryptionKeyClearsPreviousKey(t *testing.T) {
	t.Parallel()

	svc := New("device-key")
	require.NoError(t, svc.SetEncryptionKey(strings.Repeat("k", aesKeyLen256)))
	require.NoError(t, svc.SetEncryptionKey(""))
	assert.Empty(t, svc.encryptionKey)
}

func TestEncryptBytesMatchesBarkWireFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		key                string
		expectedCiphertext string
	}{
		{
			name:               "AES-128",
			key:                strings.Repeat("k", aesKeyLen128),
			expectedCiphertext: "M31imDGJqdEEszZWcWfgiMTooJniyhmvnRIIrUSrWDVi3FxcD9cl1Dyre2wQ8a5I2bhFyBrGPeMEwp1+Uw==",
		},
		{
			name:               "AES-192",
			key:                strings.Repeat("k", aesKeyLen192),
			expectedCiphertext: "QyuvztP2Yktq2khUWXw6hYrjLnhuc+kC44lyGOaVYBRJC4vsYaB6RewayvX+4ZKjUUIifRWDtIiVFq/KgA==",
		},
		{
			name:               "AES-256",
			key:                strings.Repeat("k", aesKeyLen256),
			expectedCiphertext: "a9CkBbe9d0qkskyAzKyb4pV/WjH6D1J/JLm0EOUTMcoTpJEs10gqANkkLKaOPqVXL0yo0AaxElYjmLmynw==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ciphertext, err := encryptBytes(
				[]byte(tt.key),
				"fixed-iv-123",
				[]byte(`{"body":"Encrypted weather","level":"active"}`),
			)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedCiphertext, ciphertext)
		})
	}
}

func TestSendPlaintext(t *testing.T) {
	t.Parallel()

	var got encryptedOrPlain
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/push", r.URL.Path)
		assert.Equal(t, "application/json; charset=utf-8", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.NoError(t, json.Unmarshal(body, &got))
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	svc := NewWithServers("device-key", server.URL)
	require.NoError(t, svc.Send(t.Context(), "Title", "Body"))

	assert.Equal(t, "device-key", got.DeviceKey)
	assert.Equal(t, "Title", got.Title)
	assert.Equal(t, "Body", got.Body)
	assert.Equal(t, "alarm.caf", got.Sound)
	assert.Empty(t, got.Ciphertext)
	assert.Empty(t, got.IV)
}

func TestSendPlaintextKeepsEmptyDeviceKey(t *testing.T) {
	t.Parallel()

	var raw []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		raw, err = io.ReadAll(r.Body)
		assert.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	svc := NewWithServers("", server.URL)
	require.NoError(t, svc.Send(t.Context(), "Title", "Body"))
	assert.Contains(t, string(raw), `"device_key":""`)
}

func TestSendEncryptsCompletePayloadAndRotatesIV(t *testing.T) {
	t.Parallel()

	const (
		deviceKey = "private-device"
		title     = "Private encrypted title"
		body      = "Private encrypted body"
		key       = "kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk"
	)

	firstServerGot := make(chan encryptedPostData, 1)
	secondServerGot := make(chan encryptedPostData, 1)

	firstServer := httptest.NewServer(captureEncryptedRequest(t, firstServerGot))
	t.Cleanup(firstServer.Close)
	secondServer := httptest.NewServer(captureEncryptedRequest(t, secondServerGot))
	t.Cleanup(secondServer.Close)

	svc := NewWithServers(deviceKey, firstServer.URL, secondServer.URL)
	require.NoError(t, svc.SetEncryptionKey(key))

	require.NoError(t, svc.Send(t.Context(), title, body))

	requests := []encryptedPostData{<-firstServerGot, <-secondServerGot}
	assert.NotEqual(t, requests[0].IV, requests[1].IV)
	assert.NotEqual(t, requests[0].Ciphertext, requests[1].Ciphertext)

	expectedPlaintext := notificationParams{
		Title: title,
		Body:  body,
		Sound: "alarm.caf",
	}
	for _, request := range requests {
		assert.Equal(t, deviceKey, request.DeviceKey)
		assert.Regexp(t, `^[A-Za-z0-9_-]{12}$`, request.IV)
		assert.Equal(t, expectedPlaintext, decryptPostData(t, key, request))
	}
}

func TestSendEncryptionFailurePreventsNetworkIO(t *testing.T) {
	t.Parallel()

	called := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	svc := NewWithServers("private-device", server.URL)
	svc.encryptionKey = []byte("invalid")

	err := svc.Send(t.Context(), "private title", "private body")
	require.Error(t, err)
	assert.False(t, called)
	require.EqualError(t, err, `send message to bark server "`+server.URL+`/": encrypt payload`)
	assert.NotContains(t, err.Error(), "private title")
	assert.NotContains(t, err.Error(), "private body")
}

func TestGenerateBarkIV(t *testing.T) {
	t.Parallel()

	iv, err := generateBarkIV()
	require.NoError(t, err)
	assert.Regexp(t, `^[A-Za-z0-9_-]{12}$`, iv)
}

type encryptedOrPlain struct {
	DeviceKey  string `json:"device_key"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Sound      string `json:"sound"`
	Ciphertext string `json:"ciphertext"`
	IV         string `json:"iv"`
}

func captureEncryptedRequest(t *testing.T, got chan<- encryptedPostData) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		raw, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var payload encryptedPostData
		assert.NoError(t, json.Unmarshal(raw, &payload))
		got <- payload
		w.WriteHeader(http.StatusOK)
	}
}

func decryptPostData(t *testing.T, key string, payload encryptedPostData) notificationParams {
	t.Helper()

	block, err := aes.NewCipher([]byte(key))
	require.NoError(t, err)
	gcm, err := cipher.NewGCM(block)
	require.NoError(t, err)

	raw, err := base64.StdEncoding.DecodeString(payload.Ciphertext)
	require.NoError(t, err)
	plaintext, err := gcm.Open(nil, []byte(payload.IV), raw, nil)
	require.NoError(t, err)

	var params notificationParams
	require.NoError(t, json.Unmarshal(plaintext, &params))

	var envelope map[string]any
	require.NoError(t, json.Unmarshal(plaintext, &envelope))
	_, hasDeviceKey := envelope["device_key"]
	assert.False(t, hasDeviceKey)

	return params
}
