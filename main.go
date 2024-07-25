package main

import (
	"github.com/joho/godotenv"

	"oh-my-chat/settings"
	"oh-my-chat/src/app"
	"oh-my-chat/src/config"
	"oh-my-chat/src/logger"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		panic("Error loading .env file")
	}

	logger.InitLog(settings.GetEnvOrDefault("LOGGER", "develop"))

	config := config.OhMyChatConfig{
		Worker:       config.Worker{Enabled: true, Number: 1},
		Api:          config.Api{Enabled: false},
		Connector:    config.Connector{Provider: config.Cli},
		ChatDatabase: config.ChatDatabase{Kind: "memory"},
	}

	app.Run(config)

}
