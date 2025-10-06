package centrifugo

import (
	"context"
	"testing"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

func TestService_Send(t *testing.T) {
	mock := &MockClient{
		PublishFunc: func(_ context.Context, channel string, data []byte) (centrifuge.PublishResult, error) {
			if channel != "test-channel" {
				t.Errorf("expected channel 'test-channel', got '%s'", channel)
			}
			if string(data) != "Test Subject\nHello, Centrifugo!" {
				t.Errorf("unexpected message: %s", string(data))
			}
			return centrifuge.PublishResult{}, nil
		},
	}
	svc := NewWithClient(mock, "test-channel")
	ctx := context.Background()
	subject := "Test Subject"
	msg := "Hello, Centrifugo!"
	if err := svc.Send(ctx, subject, msg); err != nil {
		t.Errorf("failed to send message: %v", err)
	}
}
