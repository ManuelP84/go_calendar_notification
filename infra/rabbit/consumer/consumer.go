package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/ManuelP84/calendar_notification/infra/rabbit"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

const (
	exchangeName     = "exchange"
	eventHandlerType = "eventHandler"
	prefetchCount    = 10
	eventQueue       = "eventQueue"
	concurrency      = 5
)

type Consumer struct {
	connection *amqp.Connection
}

func NewConsumer(settings *rabbit.RabbitSettings) *Consumer {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%s", settings.User, settings.Password, settings.Host, settings.Port)
	conn, err := amqp.Dial(addr)

	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	return &Consumer{conn}
}

func (c *Consumer) Run(ctx context.Context) error {
	ch, err := c.connection.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a channel", err)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("%s: %s", "Failed to declare an exchange", err)
	}

	g, groupCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return c.Listen(groupCtx, eventHandlerType)
	})
	return g.Wait()
}

func (c *Consumer) Listen(ctx context.Context, eventHandlerType string) error {
	ch, err := c.connection.Channel()

	if err != nil {
		log.Panicf("%s: %s", "Failed to open a channel", err)
	}

	err = ch.Qos(prefetchCount, 0, false)

	q, err := ch.QueueDeclare(
		eventQueue,
		false,
		false,
		true,
		false,
		nil,
	)

	err = ch.QueueBind(
		q.Name,
		"",
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		log.Panicf("%s: %s", "Failed to bind a queue", err)
	}

	mssgs, err := ch.Consume(
		q.Name,
		"consumer",
		false,
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
			log.Println(string(delivery.Body))
			return nil
		})
	}
	gCtx.Done()

	return nil
}
