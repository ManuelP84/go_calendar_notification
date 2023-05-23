package app

import (
	"context"
	"log"

	"github.com/ManuelP84/calendar_notification/infra/rabbit/consumer"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Settings    *AppSettings
	BusConsumer *consumer.Consumer
}

func NewApp(settings *AppSettings, busConsumer *consumer.Consumer) *App {
	return &App{
		Settings:    settings,
		BusConsumer: busConsumer,
	}
}

func (app *App) Run(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	app.BusConsumer.AddEventHandler("taskCreated", func(ctx context.Context, s string) error { log.Println("Handling task created!"); return nil })

	g.Go(func() error {
		if err := app.BusConsumer.Run(ctx); err != nil {
			return err
		}
		return nil
	})

	log.Println("Listening events...")

	return g.Wait()
}
