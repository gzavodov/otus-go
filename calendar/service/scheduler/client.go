package scheduler

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/gzavodov/otus-go/calendar/pkg/queue"
	"go.uber.org/zap"
)

//NewClient creates new scheduler client
func NewClient(ctx context.Context, channel queue.NotificationChannel, logger *zap.Logger) *Client {
	return &Client{
		ctx:             ctx,
		channel:         channel,
		shutdownChannel: make(chan struct{}),
		logger:          logger,
	}
}

//Client event reminder service
type Client struct {
	ctx             context.Context
	channel         queue.NotificationChannel
	shutdownChannel chan struct{}
	once            sync.Once
	logger          *zap.Logger
}

//GetServiceName returns service name
func (c *Client) GetServiceName() string {
	return "Scheduler client"
}

//Start starts scheduler client
func (c *Client) Start() error {
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
				c.LogError(ErrorCategoryChannel, err)
				return err
			}
			fmt.Fprintf(
				os.Stdout,
				"%s: %s\n",
				result.Notification.StartTime.Format("2 Jan 2006 15:04"),
				result.Notification.Title,
			)
		}
	}
}

//Stop stop scheduler server
func (c *Client) Stop() {
	c.once.Do(func() { c.shutdownChannel <- struct{}{} })
}

//LogError writes error in log
func (c *Client) LogError(name string, err error) {
	c.logger.Error(
		"Scheduler client",
		zap.NamedError(name, err),
	)
}
