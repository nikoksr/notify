// Package googlechat provides message notification integration sent to multiple
// spaces within a Google Chat Application.
package googlechat

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"google.golang.org/api/chat/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

//go:generate mockery --name=spacesMessageCreator --output=. --case=underscore --inpackage
type spacesMessageCreator interface {
	Create(string, *chat.Message) createCall
}

//go:generate mockery --name=createCall --output=. --case=underscore --inpackage
type createCall interface {
	Do(...googleapi.CallOption) (*chat.Message, error)
}

// Compile-time check to ensure that client implements the spaces message service
// interface.
var (
	_ spacesMessageCreator = new(messageCreator)
	// Compile-time check to ensure that client implements the create call interface.
	_ createCall = new(chat.SpacesMessagesCreateCall)
)

// messageCreator is wrapper struct for the native chat.SpacesMessagesService struct.
// This exists so that we can mock the chat.SpacesMessagesCreateCall with the common
// interface "createCall".
type messageCreator struct {
	*chat.SpacesMessagesService
}

func newMessageCreator(options ...option.ClientOption) (spacesMessageCreator, error) {
	ctx := context.Background()
	svc, err := chat.NewService(ctx, options...)
	if err != nil {
		return nil, err
	}
	return &messageCreator{svc.Spaces.Messages}, nil
}

func (m *messageCreator) Create(parent string, message *chat.Message) createCall {
	return m.SpacesMessagesService.Create(parent, message)
}

// Service encapsulates the google chat client along with internal state for storing
// chat spaces.
type Service struct {
	messageCreator spacesMessageCreator
	spaces         []string
}

// New returns an instance of the google chat notification service
func New(options ...option.ClientOption) (*Service, error) {
	svc, err := newMessageCreator(options...)
	if err != nil {
		return nil, err
	}
	s := &Service{
		messageCreator: svc,
		spaces:         []string{},
	}
	return s, nil
}

// AddReceivers takes a name of authorized spaces and appends them to the internal
// spaces slice. The Send method will send a given message to all those spaces.
func (s *Service) AddReceivers(spaces ...string) {
	s.spaces = append(s.spaces, spaces...)
}

// Send takes a message subject and a message body and sends them to all the spaces
// previously set.
func (s *Service) Send(ctx context.Context, subject, message string) error {
	// Treating subject as message title
	msg := &chat.Message{Text: subject + "\n" + message}
	for _, space := range s.spaces {
		parent := fmt.Sprintf("spaces/%s", space)
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := s.messageCreator.Create(parent, msg).Do()
			if err != nil {
				return errors.Wrapf(err, "failed to send message to the google chat space: %s", space)
			}
		}
	}
	return nil
}
