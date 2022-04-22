package notify

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
