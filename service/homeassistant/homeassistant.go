package homeassistant

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// HomeAssistant struct holds necessary data to communicate with Home Assistant API.
type HomeAssistant struct {
	serverURL   string
	accessToken string
	client      *http.Client
}

// NotificationData represents the Home Assistant notification payload.
type NotificationData struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message"`
	Target  string `json:"target,omitempty"`
}

// ServiceCall represents the Home Assistant service call structure.
type ServiceCall struct {
	ServiceData NotificationData `json:"service_data"`
}

// New returns a new instance of Home Assistant notification service.
func New(serverURL, accessToken string) *HomeAssistant {
	if !strings.HasPrefix(serverURL, "http") {
		serverURL = "https://" + serverURL
	}

	return &HomeAssistant{
		serverURL:   strings.TrimSuffix(serverURL, "/"),
		accessToken: accessToken,
		client:      &http.Client{},
	}
}

// Send takes a message subject and a message body and sends them to Home Assistant.
func (h *HomeAssistant) Send(ctx context.Context, subject, message string) error {
	if h.accessToken == "" {
		return errors.New("access token is required")
	}

	serviceCall := ServiceCall{
		ServiceData: NotificationData{
			Title:   subject,
			Message: message,
		},
	}

	payload, err := json.Marshal(serviceCall)
	if err != nil {
		return fmt.Errorf("failed to marshal service call: %w", err)
	}

	url := fmt.Sprintf("%s/api/services/notify/persistent_notification", h.serverURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+h.accessToken)

	resp, err := h.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("home assistant API returned status %d", resp.StatusCode)
	}

	return nil
}

// SetTarget sets the target for notifications (e.g., specific device or service).
func (h *HomeAssistant) SetTarget(_ string) {
	// Target configuration can be extended in future versions
}
