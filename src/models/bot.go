package models

import (
	"oh-my-chat/settings"
	"oh-my-chat/src/config"
)

type Bot struct {
	ChatConnector  MessageConnector
	TelegramConfig TelegramConfig
	IsReady        string
	CliConfig      CliConfig
}

func NewBot(conn config.Connector) *Bot {
	return &Bot{
		ChatConnector:  MessageConnector(conn.Provider),
		TelegramConfig: TelegramConfig{Token: settings.GETENV("TELEGRAM_TOKEN")},
	}
}

type TelegramConfig struct {
	Token string
}

type ChatBot struct {
	BotName    string
	Engine     string
	WorkflowID string
}

type CliConfig struct{}
