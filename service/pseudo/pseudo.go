package pseudo

// Pseudo struct represents a dummy notification service.
type Pseudo struct{}

// New returns a new instance of a Pseudo notification service. This is used internally to initialize
// notification services list and prevent nil-reference errors.
func New() *Pseudo {
	return &Pseudo{}
}

// Send basically does nothing. Just here to conform the notify.Notifier interface.
func (Pseudo) Send(string, string) error {
	return nil
}
