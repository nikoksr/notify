package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

type Discord struct {
	client     *discordgo.Session
	channelIDs []string
}

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

func (d *Discord) AddReceivers(channelIDs ...string) {
	d.channelIDs = append(d.channelIDs, channelIDs...)
}

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
