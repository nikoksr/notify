package instagram

import (
	"fmt"
	"log"

	"github.com/ahmdrz/goinsta/v2"
	"github.com/pkg/errors"
)

// Instagram struct holds necessary data to communicate with the unofficial Instagram API.
type Instagram struct {
	client    *goinsta.Instagram
	usernames []string
}

// New returns a new instance of a Instagram notification service.
func New(username, password string) (*Instagram, error) {
	client := goinsta.New(username, password)
	err := client.Login()
	if err != nil {
		return nil, err
	}

	insta := &Instagram{
		client:    client,
		usernames: []string{},
	}

	return insta, nil
}

// AddReceivers takes Instagram usernames and adds them to the internal usernames list.
// The Send method will send a given message to all those users.
func (i *Instagram) AddReceivers(usernames ...string) {
	for _, username := range usernames {
		i.usernames = append(i.usernames, username)
	}
}

// Send takes a message subject and a message body and sends them to all previously set users.
func (i Instagram) Send(subject, message string) error {
	fullMessage := subject + "\n" + message
	for _, username := range i.usernames {
		// Search finds users with from most similar to least similar usernames
		result, err := i.client.Search.User(username)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		user := result.Users[0]
		if user.Username == username {
			// Doc says to use Conversation.Send for messages after initial message.
			// But seems like Inbox.New works for further messages, and Instagram.Conversation doesn't show any conversations.
			err = i.client.Inbox.New(&user, fullMessage)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Instagram user '%s'", username)
			}
		} else {
			cause := errors.New(fmt.Sprintf("the closest username found is '%s'", user.Username))
			return errors.Wrapf(cause, "failed to find the user with username '%s'", username)
		}
	}

	return nil
}

// Logout closes the current session to the API
func (i *Instagram) Logout() error {
	return i.client.Logout()
}
