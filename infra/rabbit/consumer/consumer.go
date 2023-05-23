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
	exchangeName     = "taskExchange"
	consumerName     = "taskEventConsumer"
	eventHandlerType = "eventHandler"
	routingKey       = "taskEvents"
	prefetchCount    = 10
	queueName        = "taskQueue"
	concurrency      = 5
	exchangeType     = "direct"
)

type EventHandlerFunc func(context.Context, string) error

type Consumer struct {
	connection    *amqp.Connection
	eventHandlers map[string]EventHandlerFunc
}

func NewConsumer(settings *rabbit.RabbitSettings) *Consumer {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%s/", settings.User, settings.Password, settings.Host, settings.Port)
	conn, err := amqp.Dial(addr)

	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	return &Consumer{conn, map[string]EventHandlerFunc{}}
}

func (c *Consumer) AddEventHandler(event string, handler EventHandlerFunc) {
	c.eventHandlers[event] = handler
}

func (c *Consumer) HandleEvent(ctx context.Context, data amqp.Delivery, semaphore chan bool) error {
	defer func() {
		<-semaphore
	}()

	strData := string(data.Body)

	handler, exists := c.eventHandlers[strData]

	if !exists {
		log.Println("message without a handler")
		err := data.Nack(false, true)

		return err
	}

	err := handler(ctx, strData)

	if err != nil {
		log.Println("error handling data...")
	}

	return err
}

func (c *Consumer) Run(ctx context.Context) error {

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
			// log.Println(string(delivery.Body))
			// return nil
			return c.HandleEvent(ctx, delivery, semaphore)
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
