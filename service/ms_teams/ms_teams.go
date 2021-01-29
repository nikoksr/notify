package ms_teams

import (
	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/pkg/errors"
)

// MSTeams struct holds necessary data to communicate with the MSTeams API.
type MSTeams struct {
	client   goteamsnotify.API
	webHooks []string
}

// New returns a new instance of a MSTeams notification service.
// For more information about telegram api token:
//    -> https://github.com/atc0005/go-teams-notify#example-basic
func New(apiToken string) (*MSTeams, error) {
	client := goteamsnotify.NewClient()

	m := &MSTeams{
		client:   client,
		webHooks: []string{},
	}

	return m, nil
}

// DisableWebhookValidation disables the validation webhook URLs, including the validation of known prefixes so that
// custom/private webhook URL endpoints can be used (e.g., testing purposes).
// For more information about telegram api token:
//    -> https://github.com/atc0005/go-teams-notify#example-disable-webhook-url-prefix-validation
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
//    -> https://github.com/atc0005/go-teams-notify#example-basic
func (m MSTeams) Send(subject, message string) error {
	msgCard := goteamsnotify.NewMessageCard()
	msgCard.Title = subject
	msgCard.Text = message

	for _, webHook := range m.webHooks {
		err := m.client.Send(webHook, msgCard)
		if err != nil {
			return errors.Wrapf(err, "failed to send message to Microsoft Teams via webhook '%s'", webHook)
		}
	}

	return nil
}
