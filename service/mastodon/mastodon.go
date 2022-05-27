package mastodon

import (
	"context"
	"fmt"

	"github.com/mattn/go-mastodon"
	"github.com/pkg/errors"
)

type Mastodon struct {
	client      *mastodon.Client
	mastodonIDs []string
}

type Credentials struct {
	Server       string
	ClientID     string
	ClientSecret string
	AccessToken  string
}

func New(credentials Credentials) (*Mastodon, error) {
	config := mastodon.Config(credentials)
	client := mastodon.NewClient(&config)

	// Verify Credentials
	_, err := client.GetAccountCurrentUser(context.Background())
	if err != nil {
		return nil, err
	}

	t := &Mastodon{
		client:      client,
		mastodonIDs: []string{},
	}

	return t, nil
}

func (t *Mastodon) AddReceivers(mastodonIDs ...string) {
	t.mastodonIDs = append(t.mastodonIDs, mastodonIDs...)
}

func (t Mastodon) Send(ctx context.Context, subject, message string) error {

	for _, mastodonID := range t.mastodonIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			directMessage := &mastodon.Toot{
				Status:     fmt.Sprintf("@%s, %s: %s", mastodonID, subject, message),
				Visibility: "direct",
			}

			_, err := t.client.PostStatus(ctx, directMessage)
			if err != nil {
				return errors.Wrapf(err, "failed to send direct message to mastodon ID '%s'", mastodonID)
			}

		}
	}

	return nil
}
