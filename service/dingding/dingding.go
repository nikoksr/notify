package dingding

import (
	"context"
	"fmt"

	"github.com/blinkbean/dingtalk"
	"github.com/pkg/errors"
)

type Service struct {
	config Config
	client *dingtalk.DingTalk
}

type Config struct {
	token  string
	secret string
}

func New(cfg *Config) *Service {
	dt := dingtalk.InitDingTalkWithSecret(cfg.token, cfg.secret)
	s := Service{
		config: *cfg,
		client: dt,
	}
	return &s
}

// Send takes a message subject and a message content and sends them to all previously set users.
func (s *Service) Send(ctx context.Context, subject string, content string) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		text := fmt.Sprintf("%s\n%s", subject, content)
		err := s.client.SendTextMessage(text)
		if err != nil {
			return errors.Wrapf(err, "failed to send message")
		}
	}

	return nil
}
