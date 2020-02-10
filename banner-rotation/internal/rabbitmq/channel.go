package rabbitmq

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gzavodov/otus-go/banner-rotation/queue"
	"github.com/streadway/amqp"
)

//NewChannel creates new RabbitMQ channel
func NewChannel(ctx context.Context, name string, address string) *Channel {
	return &Channel{Name: name, Address: address, ctx: ctx, mu: sync.RWMutex{}}
}

//Channel wrapper for amqp.Channel
type Channel struct {
	Address string
	Name    string

	isOpened bool
	ctx      context.Context
	conn     *amqp.Connection
	ch       *amqp.Channel
	mu       sync.RWMutex
}

//Open opens AMQP channel
func (c *Channel) Open() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isOpened {
		return nil
	}

	conn, err := amqp.Dial(c.Address)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	c.isOpened = true
	c.conn = conn
	c.ch = ch

	return nil
}

//IsOpened сhecksf channel is opened
func (c *Channel) IsOpened() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.isOpened
}

//Close closes undelying channel and connection
func (c *Channel) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.ch.Close()
	c.conn.Close()
	c.isOpened = false
}

//Write writes notification in queue
func (c *Channel) Write(item *queue.Notification) error {
	if err := c.Open(); err != nil {
		return err
	}
	defer c.Close()

	q, err := c.ch.QueueDeclare(
		c.Name, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(*item)
	if err != nil {
		return err
	}

	return c.ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

//Read creates notification read channel
func (c *Channel) Read() (<-chan *queue.ReadResult, error) {
	if err := c.Open(); err != nil {
		return nil, err
	}

	q, err := c.ch.QueueDeclare(
		c.Name, // name
		false,  // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	if err != nil {
		return nil, err
	}

	deliveryCh, err := c.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return nil, err
	}

	notificationCh := make(chan *queue.ReadResult)
	go func() {
		defer close(notificationCh)
		defer c.Close()
	label:
		for {
			select {
			case <-c.ctx.Done():
				break label
			case message, ok := <-deliveryCh:
				if !ok {
					break label
				}

				notification := &queue.Notification{}
				err := json.Unmarshal(message.Body, notification)
				if err != nil {
					notificationCh <- &queue.ReadResult{Error: err}
					break label
				}

				notificationCh <- &queue.ReadResult{Notification: notification}
			}
		}
	}()

	return notificationCh, nil
}
