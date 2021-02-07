package notify

// Notifier defines the behavior for notification services. It implements Send and AddReciever
//
// The Send command simply sends a message string to the internal destination Notifier.
//  E.g for telegram it sends the message to the specified group chat.
//
// The AddReceivers takes one or many strings and
// adds these to the list of destinations for receiving messages
// e.g. slack channels, telegram chats, email addresses.
type Notifier interface {
	Send(string, string) error
	AddReceivers(...string)
}
