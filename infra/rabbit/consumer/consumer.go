package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/ManuelP84/calendar_notification/domain/task/events"
	"github.com/ManuelP84/calendar_notification/infra/rabbit"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

const (
	exchangeName     = "taskExchange"
	consumerName     = "taskEventConsumer"
	eventHandlerType = "eventHandler"
	routingKey       = "taskEvents"
	prefetchCount    = 10
	queueName        = "taskQueue"
	concurrency      = 5
	exchangeType     = "direct"
)

type Consumer struct {
	connection    *amqp.Connection
	eventHandlers events.EventHandlers
}

func NewConsumer(settings *rabbit.RabbitSettings) *Consumer {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%s/", settings.User, settings.Password, settings.Host, settings.Port)
	conn, err := amqp.Dial(addr)

	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	return &Consumer{conn, map[string]events.EventHandlerFunc{}}
}

func (c *Consumer) AddEventHandler(event string, handler events.EventHandlerFunc) {
	c.eventHandlers[event] = handler
}

func (c *Consumer) Run(ctx context.Context) error {

	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		//1. Gorutine to listen event types
		return c.Listen(groupCtx, eventHandlerType)
	})
	return g.Wait()
}

func (c *Consumer) Listen(ctx context.Context, eventHandlerType string) error {
	ch, err := c.connection.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a channel", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("%s: %s", "Failed to declare an exchange", err)
	}

	err = ch.Qos(prefetchCount, 0, false)

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("%s: %s", "Failed to declare an queue", err)
	}

	err = ch.QueueBind(
		q.Name,
		routingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("%s: %s", "Failed to bind a queue", err)
	}

	mssgs, err := ch.Consume(
		queueName,
		consumerName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("%s: %s", "Failed to consume the messages", err)
	}

	semaphore := make(chan bool, concurrency)

	g, gCtx := errgroup.WithContext(ctx)

	for mssge := range mssgs {
		delivery := mssge
		semaphore <- true

		g.Go(func() error {
			return rabbit.HandleEvent(ctx, delivery, semaphore, c.eventHandlers)
		})
	}
	gCtx.Done()

	if !ch.IsClosed() {
		err = ch.Cancel(consumerName, false)

		if err != nil {
			return err
		}
		err = ch.Close()
		if err != nil {
			return err
		}

	}

	return g.Wait()
}
