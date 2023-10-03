package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MessageValidationTest struct {
	msg Message
	err error
}

func TestMessageValidation(t *testing.T) {
	t.Parallel()

	msgs := []MessageValidationTest{
		{
			msg: Message{Type: "stream", To: "Sai Sumith", Topic: "something", Content: "content"},
			err: nil,
		},
		{
			msg: Message{Type: "direct", To: "Sai Sumith", Topic: "something", Content: "content"},
			err: nil,
		},
		{
			msg: Message{Type: "private", To: "Sai Sumith", Topic: "something", Content: "content"},
			err: nil,
		},
		{
			msg: Message{Type: "stream", To: "Sai Sumith", Topic: "", Content: "content"},
			err: ErrInvalidMessageTopic,
		},
		{
			msg: Message{Type: "stream", To: "Sai Sumith", Topic: "something", Content: ""},
			err: ErrInvalidMessageContent,
		},
		{
			msg: Message{Type: "stream", To: 1, Topic: "something", Content: "content"},
			err: nil,
		},
		{
			msg: Message{Type: "stream", To: []int{1, 2, 3}, Topic: "something", Content: "content"},
			err: nil,
		},
		{
			msg: Message{Type: "stream", To: []string{"Sai", "Sumith"}, Topic: "something", Content: "content"},
			err: nil,
		},
		{
			msg: Message{Type: "somethingRandom", To: []string{"Sai", "Sumith"}, Topic: "something", Content: "content"},
			err: ErrInvalidMessageType,
		},
		{
			msg: Message{Type: "stream", To: true, Topic: "something", Content: "content"},
			err: ErrInvalidMessageTo,
		},
	}

	for i, v := range msgs {
		if assert.Equal(t, v.msg.Validate(), v.err) {
			t.Logf("TEST %d: validation passed", i)
		} else {
			t.Errorf("TEST %d: validation failed", i)
		}
	}
}
