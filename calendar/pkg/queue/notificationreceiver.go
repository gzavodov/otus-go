package queue

//NotificationReceiver represents notification processing contract
type NotificationReceiver interface {
	Receive(n *Notification) error
}
