package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	app := CreateApp()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Fatal error %s", err.Error())
	}
}
