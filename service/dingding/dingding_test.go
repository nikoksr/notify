package dingding

import (
	"context"
	"testing"
)

func TestService_Send(t *testing.T) {
	type args struct {
		subject string
		content string
	}

	cfg := Config{
		Token:  "xxx",
		Secret: "xxx",
	}

	arg := args{
		subject: "title",
		content: "test content abcd",
	}
	s := New(&cfg)

	t.Run("send", func(t *testing.T) {
		if err := s.Send(context.Background(), arg.subject, arg.content); err != nil {
			t.Errorf("Send() error = %v, wantErr %v", err, true)
		}
	})
}
