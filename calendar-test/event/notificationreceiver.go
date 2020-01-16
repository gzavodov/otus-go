package event

import (
	"github.com/gzavodov/otus-go/calendar/pkg/queue"
)

//NewNotificationReceiver creates new calendar event notification client
func NewNotificationReceiver(channel chan<- *queue.Notification) *NotificationReceiver {
	return &NotificationReceiver{channel: channel}
}

//NotificationReceiver represents notification processing contract
type NotificationReceiver struct {
	notifications []*queue.Notification
	channel       chan<- *queue.Notification
}

//Receive represents implementation of queue::NotificationReceiver::Receive
func (r *NotificationReceiver) Receive(notification *queue.Notification) error {
	r.notifications = append(r.notifications, notification)

	if r.channel != nil {
		r.channel <- notification
	}
	return nil
}

//FindByTitle search notification by title
func (r *NotificationReceiver) FindByTitle(title string) *queue.Notification {
	for _, notification := range r.notifications {
		if notification.Title == title {
			return notification
		}
	}
	return nil
}
