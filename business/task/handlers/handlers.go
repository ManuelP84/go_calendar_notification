package handlers

import "github.com/ManuelP84/calendar_notification/domain/task/gateway/repository"

type TaskHandlers struct {
	StoreEvent *StoreEvent
}

func NewTaskHandlers(repo repository.TaskRepository) *TaskHandlers {
	return &TaskHandlers{
		StoreEvent: NewStoreEvent(repo),
	}
}
