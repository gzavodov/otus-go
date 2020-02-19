package queue

import (
	"context"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
)

//NewNotificationClient creates new notification client
func NewNotificationClient(ctx context.Context, channel NotificationChannel, receiver NotificationReceiver, logger *zap.Logger) *NotificationClient {
	client := &NotificationClient{
		ctx:             ctx,
		channel:         channel,
		shutdownChannel: make(chan struct{}),
		logger:          logger,
	}

	if receiver == nil {
		receiver = client
	}
	client.receiver = receiver

	return client
}

//NotificationClient simple notification listener service
type NotificationClient struct {
	ctx             context.Context
	channel         NotificationChannel
	receiver        NotificationReceiver
	shutdownChannel chan struct{}
	once            sync.Once
	logger          *zap.Logger
}

//Receive represents implementation of NotificationReceiver::Receive
func (c *NotificationClient) Receive(notification *Notification) error {
	_, err := fmt.Fprintf(
		os.Stdout,
		"%s: %s\n",
		notification.Time.Format("2 Jan 2006 15:04"),
		notification.EventType,
	)
	return err
}

//Start starts notification listener
func (c *NotificationClient) Start() error {
	readCh, err := c.channel.Read()
	if err != nil {
		return err
	}

	for {
		select {
		case <-c.ctx.Done():
			return nil
		case <-c.shutdownChannel:
			return nil
		case result, ok := <-readCh:
			if !ok {
				return nil
			}

			if result.Error != nil {
				c.LogError("Channel", err)
				return err
			}

			if c.receiver != nil {
				if err := c.receiver.Receive(result.Notification); err != nil {
					c.LogError("Receiver", err)
					return err
				}

			}
		}
	}
}

//Stop stop scheduler server
func (c *NotificationClient) Stop() {
	c.once.Do(func() { c.shutdownChannel <- struct{}{} })
}

//LogError writes error in log
func (c *NotificationClient) LogError(name string, err error) {
	if c.logger != nil {
		c.logger.Error(
			"Notification Client",
			zap.NamedError(name, err),
		)
	}
}
