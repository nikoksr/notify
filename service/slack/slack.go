package slack

import (
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// Slack struct holds necessary data to communicate with the Slack API.
type Slack struct {
	client     *slack.Client
	channelIDs []string
}

// New returns a new instance of a Slack notification service.
// For more information about slack api token:
//    -> https://pkg.go.dev/github.com/slack-go/slack#New
func New(apiToken string) (*Slack, error) {
	client := slack.New(apiToken)

	s := &Slack{
		client:     client,
		channelIDs: []string{},
	}

	return s, nil
}

// AddReceivers takes Slack channel IDs and adds them to the internal channel ID list. The Send method will send
// a given message to all those channels.
func (s *Slack) AddReceivers(channelIDs ...string) {
	s.channelIDs = append(s.channelIDs, channelIDs...)
}

// Send takes a message subject and a message body and sends them to all previously set channels.
func (s Slack) Send(subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	for _, channelID := range s.channelIDs {

		id, timestamp, err := s.client.PostMessage(
			channelID,
			slack.MsgOptionText(fullMessage, false),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to send message to Slack channel '%d' at time '%s'", id, timestamp)
		}
	}

	return nil
}
