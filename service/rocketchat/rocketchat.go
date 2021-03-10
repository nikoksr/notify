package rocketchat

import (
	"context"
	"net/url"

	"github.com/RocketChat/Rocket.Chat.Go.SDK/models"
	"github.com/RocketChat/Rocket.Chat.Go.SDK/rest"
	"github.com/pkg/errors"
)

// RocketChat struct holds necessary data to communicate with the RocketChat API.
type RocketChat struct {
	client  *rest.Client
	chNames []string
}

// New returns a new instance of a RocketChat notification service.
// serverUrl is the endpoint of server i.e "localhost" , scheme is protocol i.e "http/https"
// userID and token of the user sending the message.
func New(serverURL, scheme, userID, token string) (*RocketChat, error) {
	u := url.URL{
		Scheme: scheme,
		Host:   serverURL,
	}
	authInfo := models.UserCredentials{
		ID:    userID,
		Token: token,
	}

	c := rest.NewClient(&u, false)
	if err := c.Login(&authInfo); err != nil {
		return nil, err
	}

	rc := RocketChat{
		client:  c,
		chNames: []string{},
	}

	return &rc, nil
}

// AddReceivers takes Slack channel IDs and adds them to the internal channel ID list. The Send method will send
// a given message to all those channels.
func (r *RocketChat) AddReceivers(chatIDs ...string) {
	r.chNames = append(r.chNames, chatIDs...)
}

// Send takes a message subject and a message body and sends them to all previously set channels.
// user used for sending the message has to be a member of the channel.
// https://docs.rocket.chat/api/rest-api/methods/chat/postmessage
func (r *RocketChat) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	for _, chName := range r.chNames {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg := models.PostMessage{
				Channel: chName,
				Text:    fullMessage,
			}
			_, err := r.client.PostMessage(&msg)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to RocketChat channel '%s'", chName)
			}
		}
	}
	return nil
}
