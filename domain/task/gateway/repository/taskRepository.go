package repository

import (
	"context"
)

type TaskRepository interface {
	InsertEvent(ctx context.Context, eventType string) error
}
