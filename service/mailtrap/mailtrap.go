// Package mailtrap provides a notify.Notifier implementation that sends transactional emails through the Mailtrap Email
// Sending API.
package mailtrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultSendHost = "https://send.api.mailtrap.io"
	sendPath        = "/api/send"
)

// BodyType is used to specify the format of the message body.
type BodyType int

const (
	// HTML specifies that the message body is HTML. This is the default.
	HTML BodyType = iota
	// PlainText specifies that the message body is plain text.
	PlainText
)

// httpDoer is the subset of *http.Client used by Mailtrap. It allows callers
// to inject a custom client (see WithHTTPClient) and simplifies testing.
type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Mailtrap holds the data required to communicate with the Mailtrap Email Sending API.
type Mailtrap struct {
	apiKey        string
	host          string
	path          string
	client        httpDoer
	senderAddress string
	senderName    string
	usePlainText  bool
	receivers     []recipient
}

type recipient struct {
	Email string `json:"email"`
}

type sender struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type payload struct {
	From    sender      `json:"from"`
	To      []recipient `json:"to"`
	Subject string      `json:"subject"`
	Text    string      `json:"text,omitempty"`
	HTML    string      `json:"html,omitempty"`
}

// New returns a new instance of a Mailtrap notification service.
// You will need a Mailtrap API token, which can be created at https://mailtrap.io/api-tokens.
// See https://api-docs.mailtrap.io for the full API documentation.
func New(apiKey, senderAddress string, opts ...Option) *Mailtrap {
	m := &Mailtrap{
		apiKey:        apiKey,
		host:          defaultSendHost,
		path:          sendPath,
		client:        http.DefaultClient,
		senderAddress: senderAddress,
		receivers:     []recipient{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (m *Mailtrap) AddReceivers(addresses ...string) {
	for _, address := range addresses {
		m.receivers = append(m.receivers, recipient{Email: address})
	}
}

// BodyFormat can be used to specify the format of the message body.
// Default BodyType is HTML.
func (m *Mailtrap) BodyFormat(format BodyType) {
	switch format {
	case PlainText:
		m.usePlainText = true
	case HTML:
		m.usePlainText = false
	default:
		m.usePlainText = false
	}
}

// Send takes a message subject and a message body and sends them to all previously set receivers.
// Message body supports html as markup language by default; use BodyFormat(PlainText) to send plain text instead.
func (m Mailtrap) Send(ctx context.Context, subject, message string) error {
	body := payload{
		From:    sender{Email: m.senderAddress, Name: m.senderName},
		To:      m.receivers,
		Subject: subject,
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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.host+m.path, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Token", m.apiKey)

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("the Mailtrap endpoint returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
