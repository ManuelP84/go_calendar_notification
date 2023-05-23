package handlers

import (
	"context"
	"fmt"

	"github.com/ManuelP84/calendar_notification/domain/task/gateway/repository"
)

const (
	emptyString = ""
)

type StoreEvent struct {
	TaskRepository repository.TaskRepository
}

func NewStoreEvent(repo repository.TaskRepository) *StoreEvent {
	return &StoreEvent{repo}
}

func (handler *StoreEvent) StoreTaskEvent(ctx context.Context, event string) error {
	if event == emptyString {
		return fmt.Errorf("event can't be empty")
	}

	return handler.TaskRepository.InsertEvent(ctx, event)
}
