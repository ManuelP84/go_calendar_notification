package repository

import (
	"context"

	"github.com/ManuelP84/calendar_notification/domain/task/events"
)

type TaskRepository interface {
	InsertEvent(ctx context.Context, event events.TaskEvent) error
}
