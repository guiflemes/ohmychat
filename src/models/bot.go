package models

import "oh-my-chat/settings"

type Bot struct {
	ChatConnector  MessageConnector
	TelegramConfig TelegramConfig
	IsReady        string
}

func NewBot(conn MessageConnector) *Bot {
	return &Bot{
		ChatConnector:  conn,
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
