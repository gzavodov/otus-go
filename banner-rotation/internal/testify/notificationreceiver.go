//NewNotificationReceiver creates new calendar event notification client
package testify

import (
	"fmt"
	"sync"
	"time"

	"github.com/gzavodov/otus-go/banner-rotation/queue"
)

func NewNotificationReceiver() *NotificationReceiver {
	return &NotificationReceiver{
		channel: make(chan *queue.Notification),
		mu:      sync.RWMutex{},
	}
}

//NotificationReceiver represents notification processing contract
type NotificationReceiver struct {
	notifications []*queue.Notification
	channel       chan *queue.Notification
	mu            sync.RWMutex
}

//Receive represents implementation of queue::NotificationReceiver::Receive
func (r *NotificationReceiver) Receive(notification *queue.Notification) error {
	r.mu.Lock()
	r.notifications = append(r.notifications, notification)
	r.mu.Unlock()

	if r.channel != nil {
		r.channel <- notification
	}
	return nil
}

//Find search notification by banner ID and event type
func (r *NotificationReceiver) Find(bannerID int64, eventType string) *queue.Notification {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.notifications {
		if item.BannerID == bannerID && item.EventType == eventType {
			return item
		}
	}
	return nil
}

//Wait waits for banner event notification
func (r *NotificationReceiver) Wait(bannerID int64, eventType string, timeout int) error {
	if r.Find(bannerID, eventType) != nil {
		return nil
	}

label:
	for {
		select {
		case <-time.After(time.Duration(timeout) * time.Second):
			break label
		case <-r.channel:
			if r.Find(bannerID, eventType) == nil {
				continue
			}
			return nil
		}
	}

	return fmt.Errorf("failed to receive a '%s' event notification for banner %d", eventType, bannerID)
}
