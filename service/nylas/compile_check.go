package nylas

import "github.com/nikoksr/notify"

// Compile-time check to ensure Nylas implements notify.Notifier interface.
var _ notify.Notifier = (*Nylas)(nil)
