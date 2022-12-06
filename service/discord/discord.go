package discord

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/switchupcb/disgo"
)

//go:generate mockery --name=discordSession --output=. --case=underscore --inpackage
type discordSession interface{}

// Compile-time check to ensure that disgo.Client implements the discordSession interface.
var _ discordSession = new(disgo.Client)

// Discord struct holds necessary data to communicate with the Discord API.
type Discord struct {
	client     discordSession
	channelIDs []string
}

// New returns a new instance of a Discord notification service.
func New() *Discord {
	return &Discord{
		client:     &disgo.Client{},
		channelIDs: []string{},
	}
}

// AuthenticateWithBotToken authenticates you as a bot to Discord via the given access token.
func (d *Discord) AuthenticateWithBotToken(token string) error {
	d.client = &disgo.Client{
		Authentication: disgo.BotToken(token),
		Config:         disgo.DefaultConfig(),
	}

	return nil
}

// AuthenticateWithOAuth2Token authenticates you to Discord via the given OAUTH2 token.
func (d *Discord) AuthenticateWithOAuth2Token(token string) error {
	d.client = &disgo.Client{
		Authentication: disgo.BearerToken(token),
		Config:         disgo.DefaultConfig(),
	}

	return nil
}

// AddReceivers takes Discord channel IDs and adds them to the internal channel ID list. The Send method will send
// a given message to all those channels.
func (d *Discord) AddReceivers(channelIDs ...string) {
	d.channelIDs = append(d.channelIDs, channelIDs...)
}

// Send takes a message subject and a message body and sends them to all previously set chats.
func (d Discord) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	for _, channelID := range d.channelIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			request := &disgo.CreateMessage{
				MessageReference: nil,
				Content:          disgo.Pointer(fullMessage),
				Nonce:            nil,
				TTS:              nil,
				AllowedMentions:  nil,
				Flags:            nil,
				ChannelID:        channelID,
				Embeds:           nil,
				Components:       nil,
				StickerIDS:       nil,
				Files:            nil,
				Attachments:      nil,
			}

			// assertion required due to Discord.client mock field.
			bot, ok := d.client.(*disgo.Client)
			if !ok {
				return fmt.Errorf("mock client assertion failure: failed to send message to Discord channel '%s'", channelID)
			}

			if _, err := request.Send(bot); err != nil {
				return errors.Wrapf(err, "failed to send message to Discord channel '%s'", channelID)
			}
		}
	}

	return nil
}
