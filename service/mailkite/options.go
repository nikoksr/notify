package mailkite

import (
	"net/http"
	"strings"
)

// Option describes a functional parameter for the MailKite constructor.
type Option func(*MailKite)

// WithSenderName sets the display name used for the sender ("from") address.
// If not set, only the sender email address is sent.
func WithSenderName(name string) Option {
	return func(m *MailKite) {
		m.senderName = name
	}
}

// WithReplyTo sets the Reply-To address on every message. Replies to notifications land there instead of at the
// sender address.
func WithReplyTo(address string) Option {
	return func(m *MailKite) {
		m.replyTo = address
	}
}

// WithBaseURL overrides the MailKite API base URL. This is primarily useful for testing or for pointing the client at
// a custom endpoint. A trailing slash is ignored.
func WithBaseURL(baseURL string) Option {
	return func(m *MailKite) {
		m.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithHTTPClient sets a custom HTTP client for communicating with the MailKite API.
// This is useful for configuring timeouts, proxies, or custom transports.
func WithHTTPClient(client *http.Client) Option {
	return func(m *MailKite) {
		if client != nil {
			m.client = client
		}
	}
}
