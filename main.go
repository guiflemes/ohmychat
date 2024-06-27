package main

import (
	"log"

	"github.com/joho/godotenv"

	"oh-my-chat/src/app"
	"oh-my-chat/src/config"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := config.OhMyChatConfig{
		Worker:    config.Worker{Enabled: true, Number: 1},
		Api:       config.Api{Enabled: false},
		Connector: config.Connector{Provider: config.Cli},
	}

	app.Run(config)

	// ohmychat.RunCli()
}
