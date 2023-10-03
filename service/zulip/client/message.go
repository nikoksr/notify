package client

import (
	"reflect"
)

type Message struct {
	Type    string
	To      any
	Topic   string
	Content string
}

// Validate returns an error if the message is not well-formed.
func (m *Message) Validate() error {
	switch m.Type {
	case "direct":
	case "stream":
	case "private":
		break
	default:
		return ErrInvalidMessageType
	}

	switch reflect.TypeOf(m.To) {
	case reflect.TypeOf(""):
	case reflect.TypeOf(1):
	case reflect.SliceOf(reflect.TypeOf("")):
	case reflect.SliceOf(reflect.TypeOf(1)):
		break
	default:
		return ErrInvalidMessageTo
	}

	if m.Topic == "" {
		return ErrInvalidMessageTopic
	}

	if m.Content == "" {
		return ErrInvalidMessageContent
	}

	return nil
}
