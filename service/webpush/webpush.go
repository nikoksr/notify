package webpush

import (
	"context"
	"encoding/json"
	"log"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/pkg/errors"
)

type Service struct {
	opts        webpush.Options
	subscribers []webpush.Subscription
}

// New returns a new instance of a Web Push notification service.
// It takes an email, for subscribers to respond to push notifications,
// together with a public/private key pair to encrypt notifications.
func New(email, publicKey, privateKey string) *Service {
	return &Service{
		opts: webpush.Options{
			Subscriber:      email,
			VAPIDPrivateKey: privateKey,
			VAPIDPublicKey:  publicKey,
		},
	}
}

// SubService creates a new subservice with a given topic and urgency.
// Notifications of the same topic will overwrite eachother.
// If urgent the notification will be marked high urgency.
func (s Service) SubService(topic string, urgent bool) *Service {
	s.opts.Topic = topic
	if urgent {
		s.opts.Urgency = webpush.UrgencyHigh
	}
	return &s
}

// AddRecievers takes a json response body from a subscription, unmarshals it,
// and adds it to the list of subscribers.
func (s *Service) AddReceivers(subscribers ...string) {
	for _, sub := range subscribers {
		subscriber := webpush.Subscription{}
		err := json.Unmarshal([]byte(sub), &subscriber)
		if err != nil {
			log.Printf("Cound not add subscriber: %v\n%s", err, subscriber)
			continue
		}
		s.subscribers = append(s.subscribers, subscriber)
	}
}

func (s *Service) Send(ctx context.Context, subject, message string) error {
	msg := []byte(subject + "\n" + message)

	for _, subscriber := range s.subscribers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			resp, err := webpush.SendNotification(msg, &subscriber, &s.opts)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Web Push subscriber '%s'\nStatus: %s", subscriber, resp.Status)
			}
		}
	}

	return nil
}
