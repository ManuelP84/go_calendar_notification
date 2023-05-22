package app

import (
	"github.com/ManuelP84/calendar_notification/infra/rabbit"
	"github.com/kelseyhightower/envconfig"
)

var instance *AppSettings

type AppSettings struct {
	Rabbit *rabbit.RabbitSettings
}

func loadAppSettings() *AppSettings {
	if instance == nil {
		settings := AppSettings{}

		if err := envconfig.Process("", &settings); err != nil {
			panic(err)
		}

		instance = &settings
	}
	return instance
}

func GetAppSettings() *AppSettings {
	return loadAppSettings()
}

func GetRabbitSettings() *rabbit.RabbitSettings {
	return loadAppSettings().Rabbit
}
