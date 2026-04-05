package gmail

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/mail"
	"strings"

	"golang.org/x/oauth2"
	googleGmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Gmail struct holds necessary data to communicate with the Gmail API.
type Gmail struct {
	client     *googleGmail.Service
	senderAddr string
	recipients []string
}

// New returns a new instance of a Gmail notification service.
// tokenSource provides OAuth2 credentials for the Gmail API.
// senderAddr is the email address to send from (typically the authenticated user's email).
func New(tokenSource oauth2.TokenSource, senderAddr string) (*Gmail, error) {
	srv, err := googleGmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("create gmail service: %w", err)
	}

	return &Gmail{
		client:     srv,
		senderAddr: senderAddr,
		recipients: []string{},
	}, nil
}

// AddReceivers takes email addresses and adds them to the internal recipient list.
func (g *Gmail) AddReceivers(emails ...string) {
	g.recipients = append(g.recipients, emails...)
}

// Send takes a message subject and body and sends them to all recipients via Gmail API.
func (g Gmail) Send(ctx context.Context, subject, message string) error {
	for _, recipient := range g.recipients {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			from := mail.Address{Address: g.senderAddr}
			to := mail.Address{Address: recipient}

			header := fmt.Sprintf(
				"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
				from.String(),
				to.String(),
				subject,
				message,
			)

			raw := base64.URLEncoding.EncodeToString([]byte(header))
			// Gmail API uses URL-safe base64 without padding.
			raw = strings.TrimRight(raw, "=")

			gmailMessage := &googleGmail.Message{Raw: raw}
			_, err := g.client.Users.Messages.Send("me", gmailMessage).Context(ctx).Do()
			if err != nil {
				return fmt.Errorf("send email to %q: %w", recipient, err)
			}
		}
	}

	return nil
}
