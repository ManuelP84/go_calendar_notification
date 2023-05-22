//go:build wireinject
// +build wireinject

package main

import (
	"github.com/ManuelP84/calendar_notification/infra/app"
	"github.com/ManuelP84/calendar_notification/infra/rabbit/consumer"
	"github.com/google/wire"
)

func CreateApp() *app.App {
	wire.Build(
		app.GetAppSettings,
		app.GetRabbitSettings,
		consumer.NewConsumer,
		app.NewApp,
	)
	return new(app.App)
}
