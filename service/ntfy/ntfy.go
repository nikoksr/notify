package ntfy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Ntfy struct holds necessary data to communicate with ntfy API.
type Ntfy struct {
	serverURL string
	topics    []string
	client    *http.Client
}

// Message represents the ntfy message payload.
type Message struct {
	Topic   string `json:"topic"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message"`
	Tags    string `json:"tags,omitempty"`
}

// New returns a new instance of Ntfy notification service.
func New(serverURL string) *Ntfy {
	if !strings.HasPrefix(serverURL, "http") {
		serverURL = "https://" + serverURL
	}
	
	return &Ntfy{
		serverURL: strings.TrimSuffix(serverURL, "/"),
		client:    &http.Client{},
		topics:    make([]string, 0),
	}
}

// AddReceivers takes ntfy topic names and adds them to the internal topics list.
func (n *Ntfy) AddReceivers(topics ...string) {
	n.topics = append(n.topics, topics...)
}

// Send takes a message subject and a message body and sends them to all previously set topics.
func (n *Ntfy) Send(ctx context.Context, subject, message string) error {
	if len(n.topics) == 0 {
		return fmt.Errorf("no topics configured")
	}

	for _, topic := range n.topics {
		if err := n.sendToTopic(ctx, topic, subject, message); err != nil {
			return fmt.Errorf("failed to send to topic %s: %w", topic, err)
		}
	}

	return nil
}

func (n *Ntfy) sendToTopic(ctx context.Context, topic, subject, message string) error {
	msg := Message{
		Topic:   topic,
		Title:   subject,
		Message: message,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.serverURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ntfy API returned status %d", resp.StatusCode)
	}

	return nil
}