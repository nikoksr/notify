package bark

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const (
	aesKeyLen128 = 16
	aesKeyLen192 = 24
	aesKeyLen256 = 32
)

// Service allow you to configure Bark service.
type Service struct {
	deviceKey     string
	client        *http.Client
	serverURLs    []string
	encryptionKey []byte
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second, //nolint: mnd // 5 seconds is a reasonable timeout for a push notification
	}
}

// DefaultServerURL is the default server to use for the bark service.
const DefaultServerURL = "https://api.day.app/"

// normalizeServerURL normalizes the server URL. It prefixes it with https:// if it's not already and appends a slash
// if it's not already there. If the serverURL is empty, the DefaultServerURL is used. We're not validating the url here
// on purpose, we leave that to the http client.
func normalizeServerURL(serverURL string) string {
	if serverURL == "" {
		return DefaultServerURL
	}

	// Normalize the url
	if !strings.HasPrefix(serverURL, "http") {
		serverURL = "https://" + serverURL
	}
	if !strings.HasSuffix(serverURL, "/") {
		serverURL += "/"
	}

	return serverURL
}

// AddReceivers adds server URLs to the list of servers to use for sending messages. We call it Receivers and not
// servers because strictly speaking, the server is still receiving the message, and additionally we're following the
// naming convention of the other services.
func (s *Service) AddReceivers(serverURLs ...string) {
	for _, serverURL := range serverURLs {
		serverURL = normalizeServerURL(serverURL)
		s.serverURLs = append(s.serverURLs, serverURL)
	}
}

// NewWithServers returns a new instance of Bark service. You can use this service to send messages to bark. You can
// specify the servers to send the messages to. By default, the service will use the default server
// (https://api.day.app/) if you don't specify any servers.
func NewWithServers(deviceKey string, serverURLs ...string) *Service {
	s := &Service{
		deviceKey:  deviceKey,
		client:     defaultHTTPClient(),
		serverURLs: make([]string, 0),
	}

	if len(serverURLs) == 0 {
		serverURLs = append(serverURLs, DefaultServerURL)
	}

	// Calling service.AddReceivers() instead of directly setting the serverURLs because we want to normalize the URLs.
	s.AddReceivers(serverURLs...)

	return s
}

// New returns a new instance of Bark service. You can use this service to send messages to bark. By default, the
// service will use the default server (https://api.day.app/).
func New(deviceKey string) *Service {
	return NewWithServers(deviceKey)
}

// SetEncryptionKey configures AES-GCM encryption for notifications.
// The key must be a raw ASCII string of 16, 24, or 32 characters, matching the
// key configured in the Bark app. Each encrypted request carries a fresh
// 12-character IV, overriding the app's saved IV. An empty key disables encryption.
// An invalid key leaves the previous configuration unchanged.
// Call SetEncryptionKey before Send; it must not run concurrently with Send
// or another call to SetEncryptionKey.
func (s *Service) SetEncryptionKey(key string) error {
	if key == "" {
		s.encryptionKey = nil
		return nil
	}

	for _, r := range key {
		if r > unicode.MaxASCII {
			return errors.New("bark encryption key must contain only ASCII characters")
		}
	}

	switch len(key) {
	case aesKeyLen128, aesKeyLen192, aesKeyLen256:
	default:
		return errors.New("bark encryption key must contain exactly 16, 24, or 32 ASCII characters")
	}

	s.encryptionKey = []byte(key)

	return nil
}

// postData is the data to send to the bark server.
type postData struct {
	DeviceKey string `json:"device_key"`
	Title     string `json:"title"`
	Body      string `json:"body,omitempty"`
	Badge     int    `json:"badge,omitempty"`
	Sound     string `json:"sound,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Group     string `json:"group,omitempty"`
	URL       string `json:"pushURL,omitempty"`
}

// encryptedPostData is the Bark /push wire format used when AES-GCM is enabled.
type encryptedPostData struct {
	DeviceKey  string `json:"device_key"`
	Ciphertext string `json:"ciphertext"`
	IV         string `json:"iv"`
}

func (s *Service) send(ctx context.Context, serverURL, subject, content string) error {
	if serverURL == "" {
		return errors.New("server url is empty")
	}

	// Marshal the message to post
	message := &postData{
		DeviceKey: s.deviceKey,
		Title:     subject,
		Body:      content,
		Sound:     "alarm.caf",
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	if len(s.encryptionKey) != 0 {
		messageJSON, err = s.encryptMessage(message)
		if err != nil {
			return fmt.Errorf("encrypt payload: %w", err)
		}
	}

	pushURL := serverURL + "push"

	// Create new request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pushURL, bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response and verify success
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark returned status code %d: %s", resp.StatusCode, string(result))
	}

	return nil
}

func (s *Service) encryptMessage(message *postData) ([]byte, error) {
	// Encrypt the notification fields populated by Send; device_key stays outside.
	plaintext, err := json.Marshal(struct {
		Title string `json:"title"`
		Body  string `json:"body,omitempty"`
		Sound string `json:"sound,omitempty"`
	}{
		Title: message.Title,
		Body:  message.Body,
		Sound: message.Sound,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal plaintext message: %w", err)
	}

	// Nine random bytes encode to 12 unpadded Base64URL characters.
	// Bark uses those ASCII bytes directly as the GCM nonce, without decoding them.
	const entropyBytes = 9
	random := make([]byte, entropyBytes)
	_, _ = rand.Read(random)
	iv := base64.RawURLEncoding.EncodeToString(random)

	ciphertext, err := encryptBytes(s.encryptionKey, iv, plaintext)
	if err != nil {
		return nil, err
	}

	encrypted := &encryptedPostData{
		DeviceKey:  s.deviceKey,
		Ciphertext: ciphertext,
		IV:         iv,
	}

	messageJSON, err := json.Marshal(encrypted)
	if err != nil {
		return nil, fmt.Errorf("marshal encrypted message: %w", err)
	}

	return messageJSON, nil
}

func encryptBytes(key []byte, iv string, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("create AES-GCM cipher: %w", err)
	}

	combined := gcm.Seal(nil, []byte(iv), plaintext, nil)

	return base64.StdEncoding.EncodeToString(combined), nil
}

// Send takes a message subject and a message content and sends them to bark application.
func (s *Service) Send(ctx context.Context, subject, content string) error {
	if s.client == nil {
		return errors.New("client is nil")
	}

	for _, serverURL := range s.serverURLs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := s.send(ctx, serverURL, subject, content)
			if err != nil {
				return fmt.Errorf("send message to bark server %q: %w", serverURL, err)
			}
		}
	}

	return nil
}
