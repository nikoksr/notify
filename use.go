package notify

// useService adds a given service to the notifiers services list. If the list still contains
// a pseudo service we remove it before adding the 'real' service.
func (n *Notify) useService(service Notifier) {
	if service == nil {
		return
	}

	n.notifiers = append(n.notifiers, service)
}

// UseService adds a given service to the notifiers services list.
func (n *Notify) UseService(service Notifier) {
	n.useService(service)
}
