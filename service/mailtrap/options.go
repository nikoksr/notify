package mailtrap

import "net/http"

const (
	bulkSendHost = "https://bulk.api.mailtrap.io"
	sandboxHost  = "https://sandbox.api.mailtrap.io"
)

// Option describes a functional parameter for the Mailtrap constructor.
type Option func(*Mailtrap)

// WithSenderName sets the display name used for the sender ("from") address.
// If not set, only the sender email address is sent.
func WithSenderName(name string) Option {
	return func(m *Mailtrap) {
		m.senderName = name
	}
}

// WithBulkHost routes messages through the Mailtrap bulk sending host, which is intended for high-volume sending.
func WithBulkHost() Option {
	return func(m *Mailtrap) {
		m.host = bulkSendHost
	}
}

// WithSandbox routes messages through the Mailtrap sandbox host, capturing emails in the given test inbox instead of
// delivering them. The sandbox API requires the inbox ID as part of the request path, so it must be provided here. See
// https://api-docs.mailtrap.io for details.
func WithSandbox(inboxID string) Option {
	return func(m *Mailtrap) {
		m.host = sandboxHost
		m.path = sendPath + "/" + inboxID
	}
}

// WithHost overrides the Mailtrap API host. This is primarily useful for testing or for pointing the client at a custom
// endpoint.
func WithHost(host string) Option {
	return func(m *Mailtrap) {
		m.host = host
	}
}

// WithHTTPClient sets a custom HTTP client for communicating with the Mailtrap API.
// This is useful for configuring timeouts, proxies, or custom transports.
func WithHTTPClient(client *http.Client) Option {
	return func(m *Mailtrap) {
		if client != nil {
			m.client = client
		}
	}
}
