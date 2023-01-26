// Package mattermost provides message notification integration for mattermost.com.
package mattermost

import (
	"context"
	"io"
	"log"
	stdhttp "net/http"

	"github.com/pkg/errors"

	"github.com/nikoksr/notify/service/http"
)

//go:generate mockery --name=httpClient --output=. --case=underscore --inpackage
type httpClient interface {
	AddReceivers(wh ...*http.Webhook)
	PreSend(prefn http.PreSendHookFn)
	Send(ctx context.Context, subject, message string) error
	PostSend(postfn http.PostSendHookFn)
}

// Service encapsulates the notify httpService client and contains mattermost channel ids.
type Service struct {
	loginClient   httpClient
	messageClient httpClient
	channelIDs    []string
}

// New returns a new instance of a Mattermost notification service.
func New(url string) *Service {
	httpService := setupMsgService(url)
	return &Service{
		setupLoginService(url, httpService),
		httpService,
		[]string{},
	}
}

// LoginWithCredentials provides helper for authentication using Mattermost user/admin credentials.
func (s *Service) LoginWithCredentials(ctx context.Context, loginID, password string) error {
	// request login
	if err := s.loginClient.Send(ctx, loginID, password); err != nil {
		return errors.Wrapf(err, "failed login to Mattermost server")
	}
	return nil
}

// AddReceivers takes Mattermost channel IDs or Chat IDs and adds them to the internal channel ID list.
// The Send method will send a given message to all these channels.
func (s *Service) AddReceivers(channelIDs ...string) {
	s.channelIDs = append(s.channelIDs, channelIDs...)
}

// Send takes a message subject and a message body and send them to added channel ids.
// you will need a 'create_post' permission for your username.
// refer https://api.mattermost.com/ for more info
func (s *Service) Send(ctx context.Context, subject, message string) error {
	for _, id := range s.channelIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// create post
			if err := s.messageClient.Send(ctx, id, subject+"\n"+message); err != nil {
				return errors.Wrapf(err, "failed to send message")
			}
		}
	}
	return nil
}

// setups main message service for creating posts
func setupMsgService(url string) *http.Service {
	// create new http client for sending messages/notifications
	httpService := http.New()

	// add custom payload builder
	httpService.AddReceivers(&http.Webhook{
		URL:         url + "/api/v4/posts",
		Header:      stdhttp.Header{},
		ContentType: "application/json",
		Method:      stdhttp.MethodPost,
		BuildPayload: func(channelID, subjectAndMessage string) (payload any) {
			return map[string]string{
				"channel_id": channelID,
				"message":    subjectAndMessage,
			}
		},
	})

	// add post-send hook for error checks
	httpService.PostSend(func(req *stdhttp.Request, resp *stdhttp.Response) error {
		if resp.StatusCode != stdhttp.StatusCreated {
			b, _ := io.ReadAll(resp.Body)
			return errors.New("failed to create post with status: " + resp.Status + " body: " + string(b))
		}
		return nil
	})
	return httpService
}

// setups login service to get token
func setupLoginService(url string, msgService *http.Service) *http.Service {
	// create another new http client for login request call.
	httpService := http.New()

	// append login path for the given mattermost server with custom payload builder.
	httpService.AddReceivers(&http.Webhook{
		URL:         url + "/api/v4/users/login",
		Header:      stdhttp.Header{},
		ContentType: "application/json",
		Method:      stdhttp.MethodPost,
		BuildPayload: func(loginID, password string) (payload any) {
			return map[string]string{
				"login_id": loginID,
				"password": password,
			}
		},
	})

	// Add pre-send hook to log the request before it is sent.
	httpService.PreSend(func(req *stdhttp.Request) error {
		log.Printf("Sending login request to %s", req.URL)
		return nil
	})

	// Add post-send hook to do error checks and log the response after it is received.
	// Also extract token from response header and set it as part of pre-send hook of main http client for further requests.
	httpService.PostSend(func(req *stdhttp.Request, resp *stdhttp.Response) error {
		if resp.StatusCode != stdhttp.StatusOK {
			b, _ := io.ReadAll(resp.Body)
			return errors.New("login failed with status: " + resp.Status + " body: " + string(b))
		}
		log.Printf("Login successful for %s", resp.Request.URL)

		// get token from header
		token := resp.Header.Get("Token")
		if token == "" {
			return errors.New("received empty token")
		}

		// set token as pre-send hook
		msgService.PreSend(func(req *stdhttp.Request) error {
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		})
		return nil
	})
	return httpService
}
