package zulip

// Receiver encapsulates a receiver credentials for a direct or stream message.
type Receiver struct {
	email  string
	stream string
	topic  string
}

// Direct specifies a Zulip Direct message
func Direct(email string) *Receiver {
	return &Receiver{email: email}
}

// Stream specifies a Zulip Stream message
func Stream(stream, topic string) *Receiver {
	return &Receiver{stream: stream, topic: topic}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
	Result  string `json:"result"`
}
