//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ManuelP84/calendar_notification/app"
	"github.com/ManuelP84/calendar_notification/business/task/handlers"
	busConsumer "github.com/ManuelP84/calendar_notification/domain/gateways/bus/consumer"
	"github.com/ManuelP84/calendar_notification/domain/task/gateway/repository"
	"github.com/ManuelP84/calendar_notification/infra/config"
	"github.com/ManuelP84/calendar_notification/infra/mongo/task"
	"github.com/ManuelP84/calendar_notification/infra/rabbit/consumer"
	"github.com/google/wire"
)

func CreateApp() *app.App {
	wire.Build(
		config.GetAppSettings,
		config.GetRabbitSettings,
		config.GetMongoSettings,
		consumer.NewConsumer,
		handlers.NewTaskHandlers,
		task.NewMongoRepository,
		wire.Bind(new(busConsumer.TaskConsumer), new(*consumer.Consumer)),
		wire.Bind(new(repository.TaskRepository), new(*task.MongoRepository)),
		app.NewApp,
	)
	return new(app.App)
}
