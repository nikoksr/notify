package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type MessageValidationTest struct {
	msg Message
	err error
}

func TestMessageValidation(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

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
		assert.Equal(
			v.msg.Validate(),
			v.err,
			fmt.Sprintf("TEST %d: validation failed", i),
		)
	}
}
