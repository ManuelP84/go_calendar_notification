//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ManuelP84/calendar_notification/business/task/handlers"
	"github.com/ManuelP84/calendar_notification/domain/task/gateway/repository"
	"github.com/ManuelP84/calendar_notification/infra/app"
	"github.com/ManuelP84/calendar_notification/infra/mongo/task"
	"github.com/ManuelP84/calendar_notification/infra/rabbit/consumer"
	"github.com/google/wire"
)

func CreateApp() *app.App {
	wire.Build(
		app.GetAppSettings,
		app.GetRabbitSettings,
		app.GetMongoSettings,
		consumer.NewConsumer,
		handlers.NewTaskHandlers,
		task.NewMongoRepository,
		wire.Bind(new(repository.TaskRepository), new(*task.MongoRepository)),
		app.NewApp,
	)
	return new(app.App)
}
