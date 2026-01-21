package zulip

import (
	"context"

	z "github.com/nikoksr/notify/service/zulip/client"
)

//go:generate mockery --name=zulipClient --output=. --case=underscore --inpackage
type zulipClient interface {
	Send(msg *z.Message) (*z.Response, error)
}

// Compile-time check to ensure that zulip client send function implements the zulipClient interface
var _ zulipClient = new(z.Client)

type Receiver struct {
	_type  string
	_to    string
	_topic string
}

type Service struct {
	client zulipClient
	recv   []*Receiver
}

func New(baseURL, botEmail, apiKey string) (*Service, error) {
	client, err := z.NewClient(
		z.WithBaseURL(baseURL),
		z.WithCreds(botEmail, apiKey),
	)
	if err != nil {
		return nil, err
	}

	service := &Service{
		client,
		[]*Receiver{},
	}

	return service, nil
}

func Direct(email string) *Receiver {
	return &Receiver{
		_type:  "direct",
		_to:    email,
		_topic: "<ignored>",
	}
}

func Stream(stream, topic string) *Receiver {
	return &Receiver{
		_type:  "stream",
		_to:    stream,
		_topic: topic,
	}
}

func (s *Service) AddReceivers(recvs ...*Receiver) {
	s.recv = append(s.recv, recvs...)
}

func (s *Service) Send(ctx context.Context, subject, message string) error {
	for _, recv := range s.recv {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg := z.Message{
				Type:    recv._type,
				To:      recv._to,
				Topic:   recv._topic,
				Content: subject + " " + message,
			}

			_, err := s.client.Send(&msg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
