package notify

import (
	"github.com/nikoksr/notify/service/pseudo"
)

// useService adds a given service to the notifiers services list. If the list still contains
// a pseudo service we remove it before adding the 'real' service.
func (n *Notify) useService(service Notifier) {
	if service == nil {
		return
	}

	// Remove pseudo service in case a 'real' service will be added
	if len(n.notifiers) > 0 {
		_, isPseudo := n.notifiers[0].(*pseudo.Pseudo)
		if isPseudo {
			n.notifiers = n.notifiers[1:]
		}
	}

	n.notifiers = append(n.notifiers, service)
}

// usePseudo adds a pseudo notification service to the notifiers services list.
func (n *Notify) usePseudo() {
	n.useService(pseudo.New())
}

// UseService adds a given service to the notifiers services list.
func (n *Notify) UseService(service Notifier) {
	n.useService(service)
}
