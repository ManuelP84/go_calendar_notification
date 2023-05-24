package events

import (
	"context"

	"github.com/ManuelP84/calendar_notification/domain/task/models"
)

type TaskEvent struct {
	EventType string
	Task      *models.Task
}

type (
	EventHandlerFunc func(context.Context, TaskEvent) error
	EventHandlers    map[string]EventHandlerFunc
)
