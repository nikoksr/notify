package notify

// useService adds a given service to the notifiers services list.
func (n *Notify) useService(service Notifier) {
	if service == nil {
		return
	}
	n.notifiers = append(n.notifiers, service)
}

// UseService adds the given service(s) to the notifiers services list.
func (n *Notify) UseService(service ...Notifier) {
	for _, s := range service {
		n.useService(s)
	}
}
