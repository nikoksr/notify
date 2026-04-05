package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
)

//go:generate mockery --name=mastodonClient --output=. --case=underscore --inpackage
type mastodonClient interface {
	PostStatus(ctx context.Context, toot *mastodon.Toot) (*mastodon.Status, error)
}

// Compile-time check
var _ mastodonClient = new(mastodon.Client)

// Mastodon struct holds necessary data to communicate with the Mastodon API.
type Mastodon struct {
	client     mastodonClient
	recipients []string
}

// New returns a new instance of a Mastodon notification service.
// serverURL is the Mastodon instance URL (e.g. "https://mastodon.social").
// clientID, clientSecret, and accessToken are OAuth2 credentials.
func New(serverURL, clientID, clientSecret, accessToken string) *Mastodon {
	client := mastodon.NewClient(&mastodon.Config{
		Server:       serverURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AccessToken:  accessToken,
	})

	return &Mastodon{
		client:     client,
		recipients: []string{},
	}
}

// AddReceivers takes Mastodon usernames (e.g. "@user@instance.social") and adds them to the
// internal recipient list. The Send method will send a direct message to all recipients.
func (m *Mastodon) AddReceivers(usernames ...string) {
	m.recipients = append(m.recipients, usernames...)
}

// Send takes a message subject and a message body and sends them as a direct message
// to all previously set recipients. Each recipient is mentioned in a separate direct toot.
func (m Mastodon) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message

	for _, recipient := range m.recipients {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			toot := &mastodon.Toot{
				Status:     recipient + " " + fullMessage,
				Visibility: "direct",
			}
			_, err := m.client.PostStatus(ctx, toot)
			if err != nil {
				return fmt.Errorf("send message to Mastodon recipient %q: %w", recipient, err)
			}
		}
	}

	return nil
}
