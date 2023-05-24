package app

import (
	"context"
	"log"

	"github.com/ManuelP84/calendar_notification/business/task/handlers"
	"github.com/ManuelP84/calendar_notification/domain/gateways/bus/consumer"
	"github.com/ManuelP84/calendar_notification/domain/task/events"
	"github.com/ManuelP84/calendar_notification/infra/config"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Settings     *config.AppSettings
	BusConsumer  consumer.TaskConsumer
	TaskHandlers *handlers.TaskHandlers
}

func NewApp(settings *config.AppSettings, busConsumer consumer.TaskConsumer, taskHandlers *handlers.TaskHandlers) *App {
	return &App{
		Settings:     settings,
		BusConsumer:  busConsumer,
		TaskHandlers: taskHandlers,
	}
}

func (app *App) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	app.BusConsumer.AddEventHandler(events.TaskCreatedEvent, func(ctx context.Context, e events.TaskEvent) error {
		log.Println("Handling task created event...")
		return app.TaskHandlers.StoreEvent.StoreTaskEvent(ctx, e)
	})

	g.Go(func() error {
		if err := app.BusConsumer.Run(ctx); err != nil {
			return err
		}
		return nil
	})

	log.Println("Listening events...")

	return g.Wait()
}
