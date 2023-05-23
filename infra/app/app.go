package app

import (
	"context"
	"log"

	"github.com/ManuelP84/calendar_notification/business/task/handlers"
	"github.com/ManuelP84/calendar_notification/domain/task"
	"github.com/ManuelP84/calendar_notification/infra/rabbit/consumer"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Settings     *AppSettings
	BusConsumer  *consumer.Consumer
	TaskHandlers *handlers.TaskHandlers
}

func NewApp(settings *AppSettings, busConsumer *consumer.Consumer, taskHandlers *handlers.TaskHandlers) *App {
	return &App{
		Settings:     settings,
		BusConsumer:  busConsumer,
		TaskHandlers: taskHandlers,
	}
}

func (app *App) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	app.BusConsumer.AddEventHandler(task.TaskCreatedEvent, func(ctx context.Context, s string) error {
		log.Println("Handling task created event...")
		return app.TaskHandlers.StoreEvent.StoreTaskEvent(ctx, s)
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
