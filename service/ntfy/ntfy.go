package ntfy

import (
	"context"
	"net/http"
	"net/url"

	"github.com/AnthonyHewins/gotfy"
	"github.com/pkg/errors"
)

// Service allow you to configure Ntfy service.
type Service struct {
	message    *gotfy.Message
	publishers []gotfy.Publisher
}

// DefaultServerURL is the default server to use for the Ntfy service.
const defaultServerURL = "https://ntfy.sh"

// AddReceivers adds server URLs to the list of servers to use for sending messages.
func (s *Service) AddPublishers(publishers ...gotfy.Publisher) {
	s.publishers = append(s.publishers, publishers...)
}

// AddMessage adds a message to be used when sending.
func (s *Service) AddMessage(m gotfy.Message) {
	s.message = &m
}

// NewWithPublishers returns a new instance of Ntfy service. You can use this service to send messages to Ntfy. You can
// specify the publishers to send the messages to. By default, the service will use the default server
// (https://ntfy.sh/) if you don't specify any servers.
func NewWithPublishers(publishers ...gotfy.Publisher) (*Service, error) {
	s := &Service{}

	if publishers == nil {
		serverURL, err := url.Parse(defaultServerURL)
		if err != nil {
			return nil, err
		}
		pub, err := gotfy.NewPublisher(serverURL, http.DefaultClient)
		if err != nil {
			return nil, err
		}
		publishers = append(publishers, *pub)
	}

	s.AddPublishers(publishers...)

	return s, nil
}

// New returns a new instance of Ntfy service. You can use this service to send messages to Ntfy. By default, the
// service will use the default server (https://ntfy.sh).
func New(hostName string, httpClient *http.Client, message *gotfy.Message) (*Service, error) {
	if hostName == "" {
		hostName = defaultServerURL
	}
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	hostNameURL, err := url.Parse(hostName)
	if err != nil {
		return nil, err
	}
	gotfyPub, err := gotfy.NewPublisher(hostNameURL, httpClient)
	if err != nil {
		return nil, err
	}

	return &Service{
		publishers: []gotfy.Publisher{*gotfyPub},
		message:    message,
	}, nil
}

func (s *Service) send(ctx context.Context, publisher *gotfy.Publisher, title, content string) error {
	if publisher == nil {
		return errors.New("publisher is nil")
	}

	if s.message == nil {
		return errors.New("message is nil. Topic needs to be set in the message")
	}

	if title != "" {
		s.message.Title = title
	}

	if content != "" {
		s.message.Message = content
	}

	_, err := publisher.SendMessage(ctx, s.message)
	if err != nil {
		return errors.Wrap(err, "sending message")
	}

	return nil
}

// Send takes a message subject and a message content and sends them to Ntfy application.
func (s *Service) Send(ctx context.Context, subject, content string) error {
	if s.publishers == nil {
		return errors.New("publishers are nil")
	}

	for _, publisher := range s.publishers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := s.send(ctx, &publisher, subject, content)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Ntfy server %q", publisher.Headers.Get("Host"))
			}
		}
	}

	return nil
}
