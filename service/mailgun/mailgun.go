package mailgun

import (
	"context"
	"fmt"

	"github.com/mailgun/mailgun-go/v5"
)

type mailgunMode int

const (
	mailgunModeText mailgunMode = iota
	mailgunModeHTML
)

// Mailgun struct holds necessary data to communicate with the Mailgun API.
type Mailgun struct {
	client            *mailgun.Client
	domain            string
	senderAddress     string
	receiverAddresses []string
	mode              mailgunMode
}

// New returns a new instance of a Mailgun notification service.
// You will need a Mailgun API key and domain name.
// See https://documentation.mailgun.com/en/latest/
func New(domain, apiKey, senderAddress string, opts ...Option) *Mailgun {
	m := &Mailgun{
		client:            mailgun.NewMailgun(apiKey),
		domain:            domain,
		senderAddress:     senderAddress,
		receiverAddresses: []string{},
		mode:              mailgunModeText,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (m *Mailgun) AddReceivers(addresses ...string) {
	m.receiverAddresses = append(m.receiverAddresses, addresses...)
}

// Send takes a message subject and a message body and sends them to all previously added email receivers.
// The body is sent as plain text by default, or as HTML when the service is configured with WithHTML().
func (m Mailgun) Send(ctx context.Context, subject, message string) error {
	var mailMessage *mailgun.PlainMessage
	switch m.mode {
	case mailgunModeText:
		mailMessage = mailgun.NewMessage(m.domain, m.senderAddress, subject, message, m.receiverAddresses...)
	case mailgunModeHTML:
		mailMessage = mailgun.NewMessage(m.domain, m.senderAddress, subject, "", m.receiverAddresses...)
		mailMessage.SetHTML(message)
	default:
		return fmt.Errorf("unknown mailgun mode: %d", m.mode)
	}

	_, err := m.client.Send(ctx, mailMessage)
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}
