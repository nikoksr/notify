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
	//nolint:staticcheck
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

	// After disabling the Notifier, Send() should return silently.
	n.WithOptions(Disable)

	if err := n.Send(ctx, "subject", "message"); err != nil {
		t.Errorf("Send() of disabled Notifier returned error: %v", err)
	}

	n.WithOptions(Enable)

	// Smuggle in a nil service. This usually never happens, since UseServices filters out nil services. But, it's good
	// to test anyway.
	n.notifiers = make([]Notifier, 0)
	n.notifiers = append(n.notifiers, nil)

	if err := n.Send(ctx, "subject", "message"); err != nil {
		t.Errorf("Send() of disabled Notifier returned no error: %v", err)
	}
	if r := recover(); r != nil {
		t.Errorf("Send() with nil service panicked: %v", r)
	}
}

func TestSendMany(t *testing.T) {
	t.Parallel()

	n := New()
	if n == nil {
		t.Fatal("New() returned nil")
	}

	var services []Notifier

	for i := 0; i < 10; i++ {
		services = append(services, mail.New("", ""))
	}

	n.UseServices(services...)

	if err := n.Send(context.Background(), "subject", "message"); err == nil {
		t.Errorf("Send() invalid mail returned no error: %v", err)
	}
}
