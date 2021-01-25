package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Discord struct holds necessary data to communicate with the Discord API.
type Discord struct {
	client     *discordgo.Session
	channelIDs []string
}

// New takes a Discord API token and returns a new instance of a Discord notification service.
func New(apiToken string) (*Discord, error) {
	client, err := discordgo.New("Bot " + apiToken)
	if err != nil {
		return nil, err
	}

	d := &Discord{
		client:     client,
		channelIDs: []string{},
	}

	return d, nil
}

// AddReceivers takes Telegram channel IDs and adds them to the internal channel ID list. The Send method will send
// a given message to all those channels.
func (d *Discord) AddReceivers(channelIDs ...string) {
	d.channelIDs = append(d.channelIDs, channelIDs...)
}

// Send takes a message subject and a message body and sends them to all previously set chats.
func (d Discord) Send(subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	for _, channelID := range d.channelIDs {
		_, err := d.client.ChannelMessageSend(channelID, fullMessage)
		if err != nil {
			return errors.Wrapf(err, "failed to send message to Discord channel '%s'", channelID)
		}
	}

	return nil
}
