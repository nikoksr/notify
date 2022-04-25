package notify

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/nikoksr/notify/service/mail"
)

func TestNew(t *testing.T) {
	t.Parallel()

	n1 := New()
	if n1 == nil {
		t.Fatal("New() returned nil")
	}
	if n1.Disabled {
		t.Fatal("New() returned disabled Notifier")
	}

	n2 := NewWithOptions()
	if n2 == nil {
		t.Fatal("NewWithOptions() returned nil")
	}
	diff := cmp.Diff(n1, n2, cmp.AllowUnexported(Notify{}))
	if diff != "" {
		t.Errorf("New() and NewWithOptions() returned different Notifiers:\n%s", diff)
	}

	n3 := NewWithOptions(Disable)
	if !n3.Disabled {
		t.Error("NewWithOptions(Disable) did not disable Notifier")
	}

	n3.WithOptions(Enable)
	if n3.Disabled {
		t.Error("WithOptions(Enable) did not enable Notifier")
	}

	n3Copy := *n3
	n3.WithOptions()
	diff = cmp.Diff(n3, &n3Copy, cmp.AllowUnexported(Notify{}))
	if diff != "" {
		t.Errorf("WithOptions() altered the Notifier:\n%s", diff)
	}

	n3.WithOptions(nil)
	if r := recover(); r != nil {
		t.Errorf("WithOptions(nil) panicked: %v", r)
	}
}

func TestDefault(t *testing.T) {
	t.Parallel()

	n := Default()
	if n == nil {
		t.Fatal("Default() returned nil")
	}
	if n.Disabled {
		t.Fatal("Default() returned disabled Notifier")
	}
	// Compare addresses on purpose.
	if n != std {
		t.Error("Default() did not return the default Notifier")
	}
}

func TestNewWithServices(t *testing.T) {
	t.Parallel()

	n1 := NewWithServices()
	if n1 == nil {
		t.Fatal("NewWithServices() returned nil")
	}

	n2 := NewWithServices(nil)
	if n2 == nil {
		t.Fatal("NewWithServices(nil) returned nil")
	}
	if len(n2.notifiers) != 0 {
		t.Error("NewWithServices(nil) did not return empty Notifier")
	}

	mailService := mail.New("", "")
	n3 := NewWithServices(mailService)
	if n3 == nil {
		t.Fatal("NewWithServices(mail.New()) returned nil")
	}
	if len(n3.notifiers) != 1 {
		t.Errorf("NewWithServices(mail.New()) was expected to have 1 notifier but had %d", len(n3.notifiers))
	} else {
		diff := cmp.Diff(n3.notifiers[0], mailService, cmp.AllowUnexported(mail.Mail{}))
		if diff != "" {
			t.Errorf("NewWithServices(mail.New()) did not correctly use service:\n%s", diff)
		}
	}
}

func TestGlobal(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if err := Send(ctx, "subject", "message"); err != nil {
		t.Errorf("Send() with no receivers returned error: %v", err)
	}

	UseServices(mail.New("", ""), nil)
	if len(std.notifiers) != 1 {
		t.Errorf("UseServices(mail.New()) was expected to have 1 notifier but had %d", len(std.notifiers))
	}

	if err := Send(ctx, "subject", "message"); err == nil {
		t.Error("Send() with invalid mail returned no error")
	}
}
