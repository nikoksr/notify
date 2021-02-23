package discord

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Discord struct holds necessary data to communicate with the Discord API.
type Discord struct {
	client     *discordgo.Session
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
func (d *Discord) authenticate(credentials ...string) error {
	client, err := discordgo.New(credentials)
	if err != nil {
		return err
	}

	client.Identify.Intents = discordgo.IntentsGuildMessageTyping

	d.client = client

	return nil
}

// AuthenticateWithCredentials authenticates you to Discord via your email and password. Note that this
// is highly discouraged by Discord. Please use an authentication token.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func (d *Discord) AuthenticateWithCredentials(email, password string) error {
	return d.authenticate(email, password)
}

// AuthenticateWithCredentialsFull authenticates you to Discord via your email, password and access token.
// This is what discord recommends.
// For more info, see here: https://pkg.go.dev/github.com/bwmarrin/discordgo@v0.22.1#New
func (d *Discord) AuthenticateWithCredentialsFull(email, password, token string, isOAuthToken bool) error {
	if isOAuthToken {
		token = parseOAuthToken(token)
	} else {
		token = parseBotToken(token)
	}

	return d.authenticate(email, password, token)
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
