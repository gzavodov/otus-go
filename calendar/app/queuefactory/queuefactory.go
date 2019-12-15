package queuefactory

import (
	"context"
	"fmt"

	"github.com/gzavodov/otus-go/calendar/app/queue"
	"github.com/gzavodov/otus-go/calendar/app/rabbitmq"
)

//Queue service Type
const (
	TypeUnknown  = 0
	TypeRabbitMQ = 1
)

//CreateQueueChannel creates queue channel by type
func CreateQueueChannel(ctx context.Context, typeID int, name string, address string) (queue.NotificationChannel, error) {
	switch typeID {
	case TypeRabbitMQ:
		return rabbitmq.NewChannel(ctx, name, address), nil
	default:
		return nil, fmt.Errorf("queue type %d is not supported in current context", typeID)
	}
}
