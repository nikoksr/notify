package textmagic

import (
	"context"
	tm "github.com/textmagic/textmagic-rest-go-v2/v2"
	"strings"
)

// TMagicService allow you to configure a TextMagic SDK client.
type TMagicService struct {
	UserName           string
	APIKey             string
	destinations       []string
	TextMagicAPIClient *tm.APIClient
}

// NewTextMagicClient creates a new text magic client
func NewTextMagicClient(userName, apiKey string) *TMagicService {

	config := tm.NewConfiguration()
	client := tm.NewAPIClient(config)

	return &TMagicService{
		TextMagicAPIClient: client,
		UserName:           userName,
		APIKey:             apiKey,
	}
}

// AddReceivers adds the given destination phone numbers to the notifier.
func (s *TMagicService) AddReceivers(phoneNumbers ...string) {
	s.destinations = append(s.destinations, phoneNumbers...)
}

// Send sends a SMS via TextMagic to all previously added receivers.
func (s *TMagicService) Send(ctx context.Context, subject, message string) error {

	// put your Username and API Key from https://my.textmagic.com/online/api/rest-api/keys page.
	auth := context.WithValue(context.Background(), tm.ContextBasicAuth, tm.BasicAuth{
		UserName: s.UserName,
		Password: s.APIKey,
	})

	text := subject + "\n" + message
	_, _, err := s.TextMagicAPIClient.TextMagicApi.SendMessage(auth, tm.SendMessageInputObject{
		Text:   text,
		Phones: strings.Join(s.destinations, ","),
	})
	if err != nil {
		return err
	}

	return nil
}
