package queue

//NotificationChannel represents event queue channel contract
type NotificationChannel interface {
	Read() (<-chan *ReadResult, error)
	Write(item *Notification) error
}
