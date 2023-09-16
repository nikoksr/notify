// Package reddit implements a Reddit notifier, allowing messages to be sent to multiple recipients
package reddit

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/pkg/errors"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

//go:generate mockery --name=redditMessageClient --output=. --case=underscore --inpackage
type redditMessageClient interface {
	Send(context.Context, *reddit.SendMessageRequest) (*reddit.Response, error)
}

// Compile-time check to ensure that reddit.MessageService implements the redditMessageClient interface.
var _ redditMessageClient = new(reddit.MessageService)

// Reddit struct holds necessary data to communicate with the Reddit API.
type Reddit struct {
	client     redditMessageClient
	recipients []string
}

// New returns a new instance of a Reddit notification service.
// For more information on obtaining client credentials:
//
//	-> https://github.com/reddit-archive/reddit/wiki/OAuth2
func New(clientID, clientSecret, username, password string) (*Reddit, error) {
	// Disable HTTP2 in http client
	// Details: https://www.reddit.com/r/redditdev/comments/t8e8hc/getting_nothing_but_429_responses_when_using_go/i18yga2/
	h := http.Client{
		Transport: &http.Transport{
			TLSNextProto: map[string]func(authority string, c *tls.Conn) http.RoundTripper{},
		},
	}
	rClient, err := reddit.NewClient(
		reddit.Credentials{
			ID:       clientID,
			Secret:   clientSecret,
			Username: username,
			Password: password,
		},
		reddit.WithHTTPClient(&h),
		reddit.WithUserAgent("github.com/heilmela/notify"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate base Reddit client")
	}

	r := &Reddit{
		client:     rClient.Message,
		recipients: []string{},
	}

	return r, nil
}

// AddReceivers takes Reddit usernames and adds them to the internal recipient list. The Send method will send
// a given message to all of those users.
func (r *Reddit) AddReceivers(recipients ...string) {
	r.recipients = append(r.recipients, recipients...)
}

// Send takes a message subject and a message body and sends them to all previously set recipients.
func (r *Reddit) Send(ctx context.Context, subject, message string) error {
	for i := range r.recipients {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m := reddit.SendMessageRequest{
				To:      r.recipients[i],
				Subject: subject,
				Text:    message,
			}

			_, err := r.client.Send(ctx, &m)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Reddit recipient '%s'", r.recipients[i])
			}
		}
	}
	return nil
}
