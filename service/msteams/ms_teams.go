package msteams

import (
	"context"

	teams "github.com/atc0005/go-teams-notify/v2"
	"github.com/pkg/errors"
)

//go:generate mockery --name=teamsClient --output=. --case=underscore --inpackage
type teamsClient interface {
	SendWithContext(ctx context.Context, webhookURL string, webhookMessage teams.MessageCard) error
	SkipWebhookURLValidationOnSend(skip bool) teams.API
}

// Compile-time check to ensure that teams.Client implements the teamsClient interface.
var _ teamsClient = teams.NewClient()

// MSTeams struct holds necessary data to communicate with the MSTeams API.
type MSTeams struct {
	client   teamsClient
	webHooks []string
}

// New returns a new instance of a MSTeams notification service.
// For more information about telegram api token:
//
//	-> https://github.com/atc0005/go-teams-notify#example-basic
func New() *MSTeams {
	client := teams.NewClient()

	m := &MSTeams{
		client:   client,
		webHooks: []string{},
	}

	return m
}

// DisableWebhookValidation disables the validation of webhook URLs, including the validation of known prefixes so that
// custom/private webhook URL endpoints can be used (e.g., testing purposes).
// For more information about telegram api token:
//
//	-> https://github.com/atc0005/go-teams-notify#example-disable-webhook-url-prefix-validation
func (m *MSTeams) DisableWebhookValidation() {
	m.client.SkipWebhookURLValidationOnSend(true)
}

// AddReceivers takes MSTeams channel web-hooks and adds them to the internal web-hook list. The Send method will send
// a given message to all those chats.
func (m *MSTeams) AddReceivers(webHooks ...string) {
	m.webHooks = append(m.webHooks, webHooks...)
}

// Send accepts a subject and a message body and sends them to all previously specified channels. Message body supports
// html as markup language.
// For more information about telegram api token:
//
//	-> https://github.com/atc0005/go-teams-notify#example-basic
func (m MSTeams) Send(ctx context.Context, subject, message string) error {
	msgCard := teams.NewMessageCard()
	msgCard.Title = subject
	msgCard.Text = message

	for _, webHook := range m.webHooks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := m.client.SendWithContext(ctx, webHook, msgCard)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Microsoft Teams via webhook '%s'", webHook)
			}
		}
	}

	return nil
}
