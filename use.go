package notify

// useService adds a given service to the Notifier's services list.
func (n *Notify) useService(service Notifier) {
	if service != nil {
		n.notifiers = append(n.notifiers, service)
	}
}

// useServices adds the given service(s) to the Notifier's services list.
func (n *Notify) useServices(service ...Notifier) {
	for _, s := range service {
		n.useService(s)
	}
}

// UseServices adds the given service(s) to the Notifier's services list.
func (n *Notify) UseServices(service ...Notifier) {
	n.useServices(service...)
}

// UseServices adds the given service(s) to the Notifier's services list.
func UseServices(service ...Notifier) {
	std.useServices(service...)
}
