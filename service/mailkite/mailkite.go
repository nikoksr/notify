// Package mailkite provides a notify.Notifier implementation that sends transactional emails through the MailKite
// Send API.
package mailkite

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	defaultBaseURL = "https://api.mailkite.dev"
	sendPath       = "/v1/send"
)

// BodyType is used to specify the format of the message body.
type BodyType int

const (
	// HTML specifies that the message body is HTML. This is the default.
	HTML BodyType = iota
	// PlainText specifies that the message body is plain text.
	PlainText
)

// httpDoer is the subset of *http.Client used by MailKite. It allows callers
// to inject a custom client (see WithHTTPClient) and simplifies testing.
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// MailKite holds the data required to communicate with the MailKite Send API.
type MailKite struct {
	apiKey        string
	baseURL       string
	client        httpDoer
	senderAddress string
	senderName    string
	replyTo       string
	usePlainText  bool
	receivers     []string
}

// payload mirrors the request body of POST /v1/send. See https://mailkite.dev/docs/sending.
type payload struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text,omitempty"`
	HTML    string   `json:"html,omitempty"`
	ReplyTo string   `json:"replyTo,omitempty"`
}

// New returns a new instance of a MailKite notification service.
// You will need a MailKite API key (mk_live_…), which can be created in the dashboard at
// https://app.mailkite.dev, and a sender address on a domain verified in that account.
// See https://mailkite.dev/docs for the full API documentation.
func New(apiKey, senderAddress string, opts ...Option) *MailKite {
	m := &MailKite{
		apiKey:        apiKey,
		baseURL:       defaultBaseURL,
		client:        http.DefaultClient,
		senderAddress: senderAddress,
		receivers:     []string{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (m *MailKite) AddReceivers(addresses ...string) {
	m.receivers = append(m.receivers, addresses...)
}

// BodyFormat can be used to specify the format of the message body.
// Default BodyType is HTML.
func (m *MailKite) BodyFormat(format BodyType) {
	switch format {
	case PlainText:
		m.usePlainText = true
	case HTML:
		m.usePlainText = false
	default:
		m.usePlainText = false
	}
}

// from renders the sender as an RFC 5322 mailbox: "Name <address>" when a sender name is set, otherwise the bare
// address.
func (m MailKite) from() string {
	if m.senderName == "" {
		return m.senderAddress
	}

	return fmt.Sprintf("%s <%s>", m.senderName, m.senderAddress)
}

// Send takes a message subject and a message body and sends them to all previously set receivers.
// Message body supports html as markup language by default; use BodyFormat(PlainText) to send plain text instead.
func (m MailKite) Send(ctx context.Context, subject, message string) error {
	body := payload{
		From:    m.from(),
		To:      m.receivers,
		Subject: subject,
		ReplyTo: m.replyTo,
	}

	if m.usePlainText {
		body.Text = message
	} else {
		body.HTML = message
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.baseURL+sendPath, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.apiKey)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(
			"the MailKite endpoint returned status %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(respBody)),
		)
	}

	return nil
}
