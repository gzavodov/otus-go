package queuemonitoring

import (
	"github.com/gzavodov/otus-go/calendar/pkg/queue"
)

func NewNotificationChannelSink(channel queue.NotificationChannel, sendingCallback func(*SendingInfo)) *NotificationChannelSink {
	return &NotificationChannelSink{
		channel:         channel,
		sendingCallback: sendingCallback,
	}
}

type NotificationChannelSink struct {
	channel         queue.NotificationChannel
	sendingCallback func(*SendingInfo)
}

func (c *NotificationChannelSink) Read() (<-chan *queue.ReadResult, error) {
	return c.channel.Read()
}

func (c *NotificationChannelSink) Write(item *queue.Notification) error {
	if err := c.channel.Write(item); err != nil {
		return err
	}

	c.RegisterSending(1)
	return nil
}

func (c *NotificationChannelSink) WriteBatch(items []*queue.Notification) error {
	if err := c.channel.WriteBatch(items); err != nil {
		return err
	}

	c.RegisterSending(len(items))
	return nil
}

func (c *NotificationChannelSink) RegisterSending(count int) {
	if c.sendingCallback == nil {
		return
	}

	go func() {
		c.sendingCallback(&SendingInfo{Count: count})
	}()
}
