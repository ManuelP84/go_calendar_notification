// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/ManuelP84/calendar_notification/business/task/handlers"
	"github.com/ManuelP84/calendar_notification/infra/app"
	"github.com/ManuelP84/calendar_notification/infra/mongo/task"
	"github.com/ManuelP84/calendar_notification/infra/rabbit/consumer"
)

// Injectors from wire.go:

func CreateApp() *app.App {
	appSettings := app.GetAppSettings()
	rabbitSettings := app.GetRabbitSettings()
	consumerConsumer := consumer.NewConsumer(rabbitSettings)
	mongoDbSettings := app.GetMongoSettings()
	mongoRepository := task.NewMongoRepository(mongoDbSettings)
	taskHandlers := handlers.NewTaskHandlers(mongoRepository)
	appApp := app.NewApp(appSettings, consumerConsumer, taskHandlers)
	return appApp
}
