package rabbit

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/ManuelP84/calendar_notification/domain/task/events"
)

func Deserialize(b []byte) (events.TaskEvent, error) {
	var taskEvent events.TaskEvent
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&taskEvent)
	return taskEvent, err
}

func HandleEvent(ctx context.Context, data amqp.Delivery, semaphore chan bool, eventHandlers events.EventHandlers) error {
	defer func() {
		<-semaphore
	}()

	taskEvent, err := Deserialize(data.Body)

	if err != nil {
		log.Panicf("%s: %s", "Failed to deserialize message", err)
	}

	handler, exists := eventHandlers[taskEvent.EventType]

	if !exists {
		log.Println("message without a handler")
		err := data.Nack(false, true)

		return err
	}

	err = handler(ctx, taskEvent)

	if err != nil {
		log.Println("error handling data...")
	}

	return err
}
