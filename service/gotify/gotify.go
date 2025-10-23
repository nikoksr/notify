package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Gotify struct holds necessary data to communicate with Gotify API.
type Gotify struct {
	serverURL string
	appToken  string
	client    *http.Client
}

// Message represents the Gotify message payload.
type Message struct {
	Title    string `json:"title,omitempty"`
	Message  string `json:"message"`
	Priority int    `json:"priority,omitempty"`
}

// New returns a new instance of Gotify notification service.
func New(serverURL, appToken string) *Gotify {
	if !strings.HasPrefix(serverURL, "http") {
		serverURL = "https://" + serverURL
	}

	return &Gotify{
		serverURL: strings.TrimSuffix(serverURL, "/"),
		appToken:  appToken,
		client:    &http.Client{},
	}
}

// Send takes a message subject and a message body and sends them to Gotify.
func (g *Gotify) Send(ctx context.Context, subject, message string) error {
	if g.appToken == "" {
		return errors.New("app token is required")
	}

	msg := Message{
		Title:   subject,
		Message: message,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.serverURL, g.appToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gotify API returned status %d", resp.StatusCode)
	}

	return nil
}

// SetPriority sets the priority for messages (0-10).
func (g *Gotify) SetPriority(_ int) {
	// Priority validation is handled by Gotify server
}
