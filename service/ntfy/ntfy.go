package ntfy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Service allow you to configure Ntfy service.
type Service struct {
	client     *http.Client
	serverURLs []string
}

func defaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

// DefaultServerURL is the default server to use for the Ntfy service.
const DefaultServerURL = "https://ntfy.sh"

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
		serverURL = serverURL + "/"
	}

	return serverURL
}

// AddReceivers adds server URLs to the list of servers to use for sending messages.
func (s *Service) AddReceivers(serverURLs ...string) {
	for _, serverURL := range serverURLs {
		serverURL = normalizeServerURL(serverURL)
		s.serverURLs = append(s.serverURLs, serverURL)
	}
}

// NewWithServers returns a new instance of Ntfy service. You can use this service to send messages to Ntfy. You can
// specify the servers to send the messages to. By default, the service will use the default server
// (https://api.day.app/) if you don't specify any servers.
func NewWithServers(serverURLs ...string) *Service {
	s := &Service{
		client: defaultHTTPClient(),
	}

	if len(serverURLs) == 0 {
		serverURLs = append(serverURLs, DefaultServerURL)
	}

	// Calling service.AddReceivers() instead of directly setting the serverURLs because we want to normalize the URLs.
	s.AddReceivers(serverURLs...)

	return s
}

// New returns a new instance of Ntfy service. You can use this service to send messages to Ntfy. By default, the
// service will use the default server (https://ntfy.sh).
func New() *Service {
	return NewWithServers()
}

// postData is the data to send to the Ntfy server.
type postData struct {
	Topic    string   `json:"topic"`
	Message  string   `json:"message,omitempty"`
	Title    string   `json:"title,omitempty"`
	Tags     []string `json:"tags,omitempty"`
	Priority int      `json:"priority,omitempty"`
	Attach   string   `json:"attach,omitempty"`
	Filename string   `json:"filename,omitempty"`
	Click    string   `json:"click,omitempty"`
	Actions  []struct {
		Action string `json:"action"`
		Label  string `json:"label,omitempty"`
		URL    string `json:"url"`
	} `json:"actions,omitempty"`
}

func (s *Service) send(ctx context.Context, serverURL, topic, content string) (err error) {
	if serverURL == "" {
		return errors.New("server url is empty")
	}

	bodyByte := []byte(content)

	isJson := json.Valid(bodyByte)
	if !isJson {
		err := errors.New("Invalid JSON Body")
		return err
	}

	var bodyj postData
	if err := json.Unmarshal(bodyByte, &bodyj); err != nil {
		err := errors.Wrap(err, "Invalid PostData structure")
		return err
	}

	if bodyj.Topic == "" {
		bodyj.Topic = topic
	}

	messageJSON, err := json.Marshal(bodyj)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}

	// Create new request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, serverURL, bytes.NewBuffer(messageJSON))
	if err != nil {
		return errors.Wrap(err, "create request")
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "send request")
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response and verify success
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "read response")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ntfy returned status code %d: %s", resp.StatusCode, string(result))
	}

	return nil
}

// Send takes a message subject and a message content and sends them to Ntfy application.
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
				return errors.Wrapf(err, "failed to send message to Ntfy server %q", serverURL)
			}
		}
	}

	return nil
}
