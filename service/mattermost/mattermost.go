// Package mattermost provides message notification integration for mattermost.com.
package mattermost

import (
	"context"
	"fmt"

	"github.com/mattermost/mattermost-server/v5/model"
)

//go:generate mockery --name=mattermostClient --output=. --case=underscore --inpackage
type mattermostClient interface {
	Login(loginID string, password string) (*model.User, *model.Response)
	CreatePost(*model.Post) (*model.Post, *model.Response)
}

// Service encapsulates the mattermost client along with internal state for storing channel ids.
type Service struct {
	client     mattermostClient
	channelIDs []string
}

// New returns a new instance of a Mattermost notification service.
func New(url string) *Service {
	return &Service{
		model.NewAPIv4Client(url),
		[]string{},
	}
}

// LoginWithCredentials provides helper for authentication using Mattermost user/admin credentials.
func (s *Service) LoginWithCredentials(loginID, password string) error {
	_, res := s.client.Login(loginID, password)
	if res.Error != nil {
		return fmt.Errorf("failed to login to Mattermost server\nstatuscode: %d\nerror: %s", res.Error.StatusCode, res.Error.Message)
	}
	return nil
}

// AddReceivers takes Mattermost channel IDs or Chat IDs and adds them to the internal channel ID list.
// The Send method will send a given message to all those channels.
func (s *Service) AddReceivers(channelIDs ...string) {
	s.channelIDs = append(s.channelIDs, channelIDs...)
}

// Send takes a message subject and a message body and sends them to all previously set channels.
// you will need a 'create_post' permission for your username.
// see https://api.mattermost.com/
func (s *Service) Send(ctx context.Context, subject, message string) error {
	for _, id := range s.channelIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg := &model.Post{
				ChannelId: id,
				Message:   subject + "\n" + message,
			}
			_, res := s.client.CreatePost(msg)
			if res.Error != nil {
				return fmt.Errorf("failed to send message to Mattermost Channel/Chat '%s'\nstatuscode: %d\nerror: %s", id, res.Error.StatusCode, res.Error.Message)
			}
		}
	}
	return nil
}
