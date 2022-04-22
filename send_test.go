package notify

import (
	"context"
	"testing"

	"github.com/nikoksr/notify/service/mail"
)

func TestNotifySend(t *testing.T) {
	t.Parallel()

	n := New()
	if n == nil {
		t.Fatal("New() returned nil")
	}
	if n.Disabled {
		t.Fatal("New() returned disabled Notifier")
	}

	// Nil context
	if err := n.Send(nil, "subject", "message"); err != nil {
		t.Errorf("Send() returned error: %v", err)
	}
	if r := recover(); r != nil {
		t.Errorf("Send() panicked: %v", r)
	}

	ctx := context.Background()
	if err := n.Send(ctx, "subject", "message"); err != nil {
		t.Errorf("Send() returned error: %v", err)
	}

	// This is not meant to test the mail service, but rather the general capability of the Send() function to catch
	// errors.
	n.UseServices(mail.New("", ""))
	if err := n.Send(ctx, "subject", "message"); err == nil {
		t.Errorf("Send() invalid mail returned no error: %v", err)
	}
}
