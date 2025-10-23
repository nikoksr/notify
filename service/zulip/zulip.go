package zulip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Zulip struct holds necessary data to communicate with Zulip API.
type Zulip struct {
	serverURL string
	email     string
	apiKey    string
	client    *http.Client
}

// Message represents the Zulip message payload.
type Message struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Topic   string `json:"topic,omitempty"`
	Content string `json:"content"`
}

// New returns a new instance of Zulip notification service.
func New(serverURL, email, apiKey string) *Zulip {
	if !strings.HasPrefix(serverURL, "http") {
		serverURL = "https://" + serverURL
	}

	return &Zulip{
		serverURL: strings.TrimSuffix(serverURL, "/"),
		email:     email,
		apiKey:    apiKey,
		client:    &http.Client{},
	}
}

// SendToStream sends a message to a Zulip stream.
func (z *Zulip) SendToStream(ctx context.Context, stream, topic, content string) error {
	if z.email == "" || z.apiKey == "" {
		return errors.New("email and API key are required")
	}

	msg := Message{
		Type:    "stream",
		To:      stream,
		Topic:   topic,
		Content: content,
	}

	return z.sendMessage(ctx, msg)
}

// SendDirectMessage sends a direct message to specific users.
func (z *Zulip) SendDirectMessage(ctx context.Context, users, content string) error {
	if z.email == "" || z.apiKey == "" {
		return errors.New("email and API key are required")
	}

	msg := Message{
		Type:    "private",
		To:      users,
		Content: content,
	}

	return z.sendMessage(ctx, msg)
}

// Send implements the notify interface - sends to a default stream.
func (z *Zulip) Send(ctx context.Context, subject, message string) error {
	return z.SendToStream(ctx, "general", subject, message)
}

func (z *Zulip) sendMessage(ctx context.Context, msg Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/messages", z.serverURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(z.email, z.apiKey)

	resp, err := z.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zulip API returned status %d", resp.StatusCode)
	}

	return nil
}
