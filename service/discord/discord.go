package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//go:generate mockery --name=discordSession --output=. --case=underscore --inpackage
type discordSession interface {
	ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
}

// Compile-time check to ensure that discordgo.Session implements the discordSession interface.
var _ discordSession = new(discordgo.Session)

// Discord struct holds necessary data to communicate with the Discord API.
type Discord struct {
	client     discordSession
	channelIDs []string
}

// New returns a new instance of a Discord notification service.
func New() *Discord {
	return &Discord{
		client:     &discordgo.Session{},
		channelIDs: []string{},
	}
}

// authenticate will try and authenticate to discord.
func (d *Discord) authenticate(token string) error {
	client, err := discordgo.New(token)
	if err != nil {
		return err
	}

	client.Identify.Intents = discordgo.IntentsGuildMessageTyping

	d.client = client

	return nil
}

// AuthenticateWithBotToken authenticates you as a bot to Discord via the given access token.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func (d *Discord) AuthenticateWithBotToken(token string) error {
	token = parseBotToken(token)

	return d.authenticate(token)
}

// AuthenticateWithOAuth2Token authenticates you to Discord via the given OAUTH2 token.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func (d *Discord) AuthenticateWithOAuth2Token(token string) error {
	token = parseOAuthToken(token)

	return d.authenticate(token)
}

// parseBotToken parses a regular token to a bot token that is understandable for discord.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func parseBotToken(token string) string {
	return "Bot " + token
}

// parseBotToken parses a regular token to a OAUTH2 token that is understandable for discord.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func parseOAuthToken(token string) string {
	return "Bearer " + token
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
			_, err := d.client.ChannelMessageSend(channelID, fullMessage)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Discord channel '%s'", channelID)
			}
		}
	}

	return nil
}
