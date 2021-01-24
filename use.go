package notify

import (
	"github.com/nikoksr/notify/service/pseudo"
)

func (n *Notifier) useService(service Service) {
	if service == nil {
		return
	}

	// Remove pseudo service in case a 'real' service will be added
	if len(n.services) > 0 {
		_, isPseudo := n.services[0].(*pseudo.Pseudo)
		if isPseudo {
			n.services = n.services[1:]
		}
	}

	n.services = append(n.services, service)
}

// usePseudo adds a pseudo Service to the Service list.
func (n *Notifier) usePseudo() {
	n.useService(pseudo.New())
}

func (n *Notifier) UseService(service Service) {
	n.useService(service)
}

/*
func (n *Notifier) UseTelegram(apiToken string, chatID int64) error {
	telegramService, err := telegram.New(apiToken)
	if err != nil {
		return err
	}

	telegramService.AddReceivers(chatID)

	n.useService(telegramService)

	return nil
}

func (n *Notifier) UseDiscordService(apiToken, channelID string) error {
	discordService, err := discord.New(apiToken)
	if err != nil {
		return err
	}

	discordService.AddReceivers(channelID)

	n.useService(discordService)

	return nil
}
*/
