package nylas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Nylas struct holds necessary data to communicate with the Nylas API v3.
type Nylas struct {
	client            httpClient
	apiKey            string
	grantID           string
	baseURL           string
	senderAddress     string
	senderName        string
	receiverAddresses []string
	usePlainText      bool
}

// httpClient interface for making HTTP requests (allows mocking in tests).
type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// BodyType is used to specify the format of the body.
type BodyType int

const (
	// PlainText is used to specify that the body is plain text.
	PlainText BodyType = iota
	// HTML is used to specify that the body is HTML.
	HTML
)

// Region represents a Nylas API region.
type Region string

const (
	// RegionUS represents the United States region.
	RegionUS Region = "US"
	// RegionEU represents the European Union region.
	RegionEU Region = "EU"
)

const (
	// BaseURLUS is the Nylas API v3 base URL for the US region.
	BaseURLUS = "https://api.us.nylas.com"
	// BaseURLEU is the Nylas API v3 base URL for the EU region.
	BaseURLEU = "https://api.eu.nylas.com"
	// DefaultBaseURL is the default Nylas API v3 base URL (US region).
	DefaultBaseURL = BaseURLUS
	// DefaultTimeout is the recommended timeout for Nylas API requests (150 seconds).
	DefaultTimeout = 150 * time.Second
)

// emailAddress represents an email recipient or sender.
type emailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

// sendMessageRequest represents the request body for sending a message via Nylas API v3.
type sendMessageRequest struct {
	To      []emailAddress `json:"to"`
	Subject string         `json:"subject"`
	Body    string         `json:"body"`
	From    []emailAddress `json:"from,omitempty"`
}

// errorResponse represents an error response from the Nylas API.
type errorResponse struct {
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
	RequestID string `json:"request_id"`
}

// New returns a new instance of a Nylas notification service for API v3.
// You will need a Nylas API key and a Grant ID.
//
// Parameters:
//   - apiKey: Your Nylas API key for authentication
//   - grantID: The Grant ID for the email account you want to send from
//   - senderAddress: The email address to send from
//   - senderName: The display name for the sender (optional, can be empty)
//
// By default, this uses the US region. For other regions, use NewWithRegion().
//
// See https://developer.nylas.com/docs/v3/getting-started/ for more information.
func New(apiKey, grantID, senderAddress, senderName string) *Nylas {
	return &Nylas{
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
		apiKey:            apiKey,
		grantID:           grantID,
		baseURL:           DefaultBaseURL,
		senderAddress:     senderAddress,
		senderName:        senderName,
		receiverAddresses: []string{},
		usePlainText:      false,
	}
}

// NewWithRegion returns a new instance of a Nylas notification service for API v3
// configured for a specific region.
//
// Parameters:
//   - apiKey: Your Nylas API key for authentication
//   - grantID: The Grant ID for the email account you want to send from
//   - senderAddress: The email address to send from
//   - senderName: The display name for the sender (optional, can be empty)
//   - region: The Nylas API region (RegionUS or RegionEU)
//
// Example:
//
//	nylasService := nylas.NewWithRegion(apiKey, grantID, email, name, nylas.RegionEU)
//
// See https://developer.nylas.com/docs/v3/getting-started/ for more information.
func NewWithRegion(apiKey, grantID, senderAddress, senderName string, region Region) *Nylas {
	n := New(apiKey, grantID, senderAddress, senderName)

	switch region {
	case RegionEU:
		n.baseURL = BaseURLEU
	case RegionUS:
		n.baseURL = BaseURLUS
	default:
		// Default to US region if unknown region is provided
		n.baseURL = BaseURLUS
	}

	return n
}

// WithBaseURL allows setting a custom base URL (e.g., for EU region: https://api.eu.nylas.com).
// This is useful for regions outside the US or for testing purposes.
func (n *Nylas) WithBaseURL(baseURL string) *Nylas {
	n.baseURL = baseURL
	return n
}

// WithHTTPClient allows setting a custom HTTP client.
// This is useful for testing or customizing timeout/transport settings.
func (n *Nylas) WithHTTPClient(client httpClient) *Nylas {
	n.client = client
	return n
}

// AddReceivers takes email addresses and adds them to the internal address list.
// The Send method will send a given message to all those addresses.
func (n *Nylas) AddReceivers(addresses ...string) {
	n.receiverAddresses = append(n.receiverAddresses, addresses...)
}

// BodyFormat can be used to specify the format of the body.
// Default BodyType is HTML.
func (n *Nylas) BodyFormat(format BodyType) {
	switch format {
	case PlainText:
		n.usePlainText = true
	case HTML:
		n.usePlainText = false
	default:
		n.usePlainText = false
	}
}

// Send takes a message subject and a message body and sends them to all previously set receivers.
// The message body supports HTML by default (unless PlainText format is specified).
//
// Note: Nylas v3 send operations are synchronous and can take up to 150 seconds for
// self-hosted Exchange servers. The timeout is set accordingly.
func (n Nylas) Send(ctx context.Context, subject, message string) error {
	if len(n.receiverAddresses) == 0 {
		return errors.New("no receivers configured")
	}

	// Build the request payload
	recipients := make([]emailAddress, 0, len(n.receiverAddresses))
	for _, addr := range n.receiverAddresses {
		recipients = append(recipients, emailAddress{
			Email: addr,
		})
	}

	body := message
	if n.usePlainText {
		// For plain text, we still send as HTML but without HTML tags
		// Nylas v3 primarily works with HTML content
		body = message
	}

	reqBody := sendMessageRequest{
		To:      recipients,
		Subject: subject,
		Body:    body,
	}

	// Add sender information if provided
	if n.senderAddress != "" {
		reqBody.From = []emailAddress{
			{
				Email: n.senderAddress,
				Name:  n.senderName,
			},
		}
	}

	// Marshal the request body
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	// Build the request URL
	url := fmt.Sprintf("%s/v3/grants/%s/messages/send", n.baseURL, n.grantID)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+n.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	// Check for success (2xx status codes)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Successfully sent
		return nil
	}

	// Handle error responses
	var errResp errorResponse
	if unmarshalErr := json.Unmarshal(respBody, &errResp); unmarshalErr != nil {
		// If we can't parse the error response, return a generic error
		return fmt.Errorf("nylas api error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return fmt.Errorf("nylas api error: %s (type: %s, request_id: %s)",
		errResp.Error.Message, errResp.Error.Type, errResp.RequestID)
}
