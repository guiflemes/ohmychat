package models

import (
	"oh-my-chat/settings"
	"oh-my-chat/src/config"
)

//TODO move Bot to config

type Bot struct {
	ChatConnector   MessageConnector
	TelegramConfig  TelegramConfig
	IsReady         string
	CliDependencies CliDependencies
}

type CliDependencies struct {
	ListWorkflows func() []string
}

func NewBot(config config.OhMyChatConfig) *Bot {
	return &Bot{
		ChatConnector:  MessageConnector(config.Connector.Provider),
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

type ChatBotCollection struct {
	bots  []*ChatBot
	names []string
}

func NewChatBotCollection(capacity int) *ChatBotCollection {
	return &ChatBotCollection{
		bots:  make([]*ChatBot, 0, capacity),
		names: make([]string, 0),
	}
}

func (c *ChatBotCollection) Add(bot *ChatBot) {
	c.bots = append(c.bots, bot)
	c.names = append(c.names, bot.BotName)
}

func (c *ChatBotCollection) Names() []string {
	return c.names
}

func (c *ChatBotCollection) Items() []*ChatBot {
	return c.bots
}
