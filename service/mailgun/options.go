package mailgun

import (
	"github.com/mailgun/mailgun-go/v5"
)

// Option describes a functional parameter for the Mailgun constructor.
type Option func(*Mailgun)

// WithEurope sets the API Mailgun base url to Europe region.
func WithEurope() Option {
	return func(m *Mailgun) {
		_ = m.client.SetAPIBase(mailgun.APIBaseEU)
	}
}

// WithHTML sets the message mode to HTML. By default, the message mode is text.
func WithHTML() Option {
	return func(m *Mailgun) {
		m.mode = mailgunModeHTML
	}
}
