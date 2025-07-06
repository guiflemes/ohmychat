package bot

import (
	"context"
	"oh-my-chat/src/connector"
)

type Bot struct {
	Connector       connector.Connector
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

func NewBot() *Bot {
	ctx, cancel := context.WithCancel(context.Background())
	return &Bot{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (b *Bot) Ctx() context.Context {
	return b.ctx
}

func (b *Bot) Shutdown() {
	b.cancel()
}
