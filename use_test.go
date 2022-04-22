package notify

import (
	"testing"

	"github.com/nikoksr/notify/service/mail"
)

func TestUseServices(t *testing.T) {
	t.Parallel()

	n := New()
	if n == nil {
		t.Fatal("New() returned nil")
	}
	if len(n.notifiers) != 0 {
		t.Fatalf("Expected len(n.notifiers) == 0, got %d", len(n.notifiers))
	}

	n.UseServices(mail.New("", ""))

	if len(n.notifiers) != 1 {
		t.Errorf("Expected len(n.notifiers) == 1, got %d", len(n.notifiers))
	}

	n.UseServices(
		mail.New("", ""),
		mail.New("", ""),
	)

	if len(n.notifiers) != 3 {
		t.Errorf("Expected len(n.notifiers) == 3, got %d", len(n.notifiers))
	}

	n.UseServices(nil)
	if r := recover(); r != nil {
		t.Errorf("Expected no panic, got %v", r)
	}
}
