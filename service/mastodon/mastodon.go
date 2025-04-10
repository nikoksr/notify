package mastodon

import (
	"context"
	"errors"

	"github.com/mattn/go-mastodon"
)

var ErrInvalidConfParams = errors.New("server and accessToken params are required")

const (
	VisibilityPublic   = "public"
	VisibilityUnlisted = "unlisted"
	VisibilityPrivate  = "private"
	VisibilityDirect   = "direct"
)

//go:generate mockery --name=mastodonClient --output=. --case=underscore --inpackage
type mastodonClient interface {
	PostStatus(ctx context.Context, toot *mastodon.Toot) (*mastodon.Status, error)
}

type Mastodon struct {
	accessToken string
	visibility  string
	client      mastodonClient
}

type Config struct {
	Visibility  string
	ServerURL   string
	AccessToken string
}

// New creates a new Mastodon client using the provided configuration.
// It returns an error if either the ServerURL or AccessToken is missing.
func New(cfg *Config) (*Mastodon, error) {
	if cfg.ServerURL == "" || cfg.AccessToken == "" {
		return nil, ErrInvalidConfParams
	}

	clientCfg := &mastodon.Config{
		Server:      cfg.ServerURL,
		AccessToken: cfg.AccessToken,
	}

	return &Mastodon{
		accessToken: cfg.AccessToken,
		visibility:  cfg.Visibility,
		client:      mastodon.NewClient(clientCfg),
	}, nil
}

// ChangeVisibility updates the default visibility setting for all future posts.
// Supported values include: public, unlisted, private (followers-only), and direct.
func (m *Mastodon) ChangeVisibility(visibility string) {
	m.visibility = visibility
}

// Send posts a new status (toot) to the configured Mastodon server.
// It uses the current visibility setting and accepts the message content as input.
func (m *Mastodon) Send(ctx context.Context, _, message string) error {
	status := &mastodon.Toot{
		Visibility: m.visibility,
		Status:     message,
	}

	if _, err := m.client.PostStatus(ctx, status); err != nil {
		return err
	}

	return nil
}
