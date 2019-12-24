package queue

//ReadResult calendar event notification read result
type ReadResult struct {
	Notification *Notification
	Error        error
}
