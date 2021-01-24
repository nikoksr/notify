package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

const (
	defaultParseMode = tgbotapi.ModeHTML
)

type Telegram struct {
	client   *tgbotapi.BotAPI
	listener *tgbotapi.BotAPI
	chatIDs  []int64
}

func New(apiToken string) (*Telegram, error) {
	client, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}

	t := &Telegram{
		client:  client,
		chatIDs: []int64{},
	}

	return t, nil
}

func (t *Telegram) AddReceivers(chatIDs ...int64) {
	t.chatIDs = append(t.chatIDs, chatIDs...)
}

func (t Telegram) Send(subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	msg := tgbotapi.NewMessage(0, fullMessage)
	msg.ParseMode = defaultParseMode

	for _, chatID := range t.chatIDs {
		msg.ChatID = chatID
		_, err := t.client.Send(msg)
		if err != nil {
			return errors.Wrapf(err, "failed to send message to Telegram chat '%d'", chatID)
		}
	}

	return nil
}
