package client

import "time"

// Hook to add baseURL to the Zulip client
func WithBaseURL(baseURL string) Option {
	return func(c *Client) error {
		c.baseURL = baseURL
		return nil
	}
}

// Hook to add necessary creds to the Zulip client
func WithCreds(email string, apiKey string) Option {
	return func(c *Client) error {
		c.email = email
		c.apiKey = apiKey
		return nil
	}
}

// Hook to add timeout to the Zulip client
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) error {
		c.timeout = timeout
		return nil
	}
}
