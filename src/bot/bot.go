package bot

import (
	"context"
	"oh-my-chat/settings"
	"oh-my-chat/src/config"
	"oh-my-chat/src/message"
)

//TODO move Bot to config

type Bot struct {
	ChatConnector   message.MessageConnector
	TelegramConfig  TelegramConfig
	IsReady         string
	CliDependencies CliDependencies
	ctx             context.Context
	cancel          context.CancelFunc
}

// CliDependencies contains the dependencies for the CliBot, including a function to list workflows
// and a flag to control the initialization of the shell.
//
// DisableInitialization is a flag that should be used exclusively during testing to prevent
// the execution of initialization code and display of messages that are specific to the production
// environment. When set to true, the CliBot will skip the usual initialization and welcome messages
// that would normally be shown during standard execution. In production environments, this flag
// should remain false to ensure that full initialization and welcome messages are displayed as expected.
type CliDependencies struct {
	ListWorkflows         func() []string
	DisableInitialization bool
}

func NewBot(config config.OhMyChatConfig) *Bot {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bot{
		ChatConnector:  message.MessageConnector(config.Connector.Provider),
		TelegramConfig: TelegramConfig{Token: settings.GetEnvOrDefault("TELEGRAM_TOKEN", "")},
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (b *Bot) Ctx() context.Context {
	return b.ctx
}

func (b *Bot) Shutdown() {
	b.cancel()
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
