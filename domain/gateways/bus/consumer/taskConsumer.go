package consumer

import (
	"context"

	"github.com/ManuelP84/calendar_notification/domain/task/events"
)

type TaskConsumer interface {
	Run(ctx context.Context) error
	Listen(ctx context.Context, eventHandlerType string) error
	AddEventHandler(event string, handler events.EventHandlerFunc)
}
