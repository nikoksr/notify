package centrifugo

import (
	"context"
	"fmt"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

// Service represents a Centrifugo notification service.
type publisher interface {
	Publish(ctx context.Context, channel string, data []byte) (centrifuge.PublishResult, error)
	Close()
}

type Service struct {
	client  publisher
	channel string
}

// New creates a new Centrifugo notification service.
// url: Centrifugo WebSocket endpoint (e.g., ws://localhost:8000/connection/websocket)
// channel: Channel to publish messages to.
// token: Optional JWT for authentication (empty string if not used).
func New(url, channel, token string) (*Service, error) {
	cfg := centrifuge.Config{}
	if token != "" {
		cfg.Token = token
	}
	client := centrifuge.NewJsonClient(url, cfg)
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("centrifugo connect error: %w", err)
	}
	return &Service{client: client, channel: channel}, nil
}

// NewWithClient allows injecting a mock client for testing.
func NewWithClient(client publisher, channel string) *Service {
	return &Service{client: client, channel: channel}
}

// Send sends a subject and message to the Centrifugo channel.
// The subject and message are concatenated with a newline.
func (s *Service) Send(ctx context.Context, subject, message string) error {
	fullMsg := subject
	if subject != "" && message != "" {
		fullMsg += "\n"
	}
	fullMsg += message
	_, err := s.client.Publish(ctx, s.channel, []byte(fullMsg))
	return err
}

// Close closes the Centrifugo client connection.
func (s *Service) Close() {
	s.client.Close()
}
