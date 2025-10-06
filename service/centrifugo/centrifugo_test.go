package centrifugo

import (
	"context"
	"testing"
)

func TestService_Send(t *testing.T) {
	// NOTE: This test requires a running Centrifugo server and a valid channel.
	// Adjust the URL, channel, and token as needed for your environment.
	url := "ws://localhost:8000/connection/websocket"
	channel := "test-channel"
	token := "" // Set JWT if required

	svc, err := New(url, channel, token)
	if err != nil {
		t.Fatalf("failed to create centrifugo service: %v", err)
	}
	defer svc.Close()

	ctx := context.Background()
	subject := "Test Subject"
	msg := "Hello, Centrifugo!"
	if err := svc.Send(ctx, subject, msg); err != nil {
		t.Errorf("failed to send message: %v", err)
	}
}
