package scheduler

import (
	"context"
	"fmt"
	"os"

	"github.com/gzavodov/otus-go/calendar/app/queue"
	"go.uber.org/zap"
)

//NewClient creates new scheduler client
func NewClient(ctx context.Context, channel queue.NotificationChannel, logger *zap.Logger) *Client {
	return &Client{ctx: ctx, channel: channel, logger: logger}
}

//Client event reminder service
type Client struct {
	ctx     context.Context
	channel queue.NotificationChannel
	logger  *zap.Logger
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

//LogError writes error in log
func (c *Client) LogError(name string, err error) {
	c.logger.Error(
		"Scheduler client",
		zap.NamedError(name, err),
	)
}
