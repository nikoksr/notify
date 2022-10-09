package mastodon

import (
    "context"
    "fmt"
    "github.com/mattn/go-mastodon"
)

// Mastodon struct holds necessary data to communicate with the Mastodon API.
type Mastodon struct {
    App    *mastodon.Application
    Client *mastodon.Client
}

type MastodonConfig struct {
    ClientID     string
    ClientSecret string
    Server       string
    ClientName   string
    Username     string
    Password     string
}

// New returns a new instance of a Mastodon notification service.
func New(config MastodonConfig) (*Mastodon, error) {
    m := Mastodon{}
    if config.ClientName == "" {
        config.ClientName = "Notify-Client"
    }

    if config.Server == "" {
        return nil, fmt.Errorf("empty server url")
    }

    if config.ClientSecret == "" || config.ClientID == "" {
        app, err := mastodon.RegisterApp(context.Background(), &mastodon.AppConfig{
            Server:     config.Server,
            ClientName: config.ClientName,
            Scopes:     "read write follow",
        })
        if err != nil {
            return nil, err
        }
        config.ClientID = app.ClientID
        config.ClientSecret = app.ClientSecret
    }
    client := mastodon.NewClient(&mastodon.Config{
        Server:       config.Server,
        ClientID:     config.ClientID,
        ClientSecret: config.ClientSecret,
    })
    err := client.Authenticate(context.Background(), config.Username, config.Password)
    if err != nil {
        return nil, err
    }
    m.Client = client
    return &m, nil
}

// Send takes a message subject and a message body and sends it to the Mastodon server.
// Plain text / emoticons are the only valid content
func (m Mastodon) Send(ctx context.Context, subject, message string) error {
    _, err := m.Client.PostStatus(ctx, &mastodon.Toot{
        Status: fmt.Sprintf("%s\n%s", subject, message),
    })
    return err
}
