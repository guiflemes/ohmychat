package core

import (
	"context"

	"go.uber.org/zap"

	"oh-my-chat/src/bot"
	"oh-my-chat/src/connector"
	"oh-my-chat/src/connector/cli"
	"oh-my-chat/src/connector/telegram"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/message"
)

type NotSupportConnectorError struct{}

func (*NotSupportConnectorError) Error() string {
	return "chat connect not supported"
}

type NewConnector func(bot *bot.Bot) (connector.Connector, error)
type GetConnectors func() map[message.MessageConnector]NewConnector

func Connectors() map[message.MessageConnector]NewConnector {
	return map[message.MessageConnector]NewConnector{
		message.Telegram: telegram.NewTelegramConnector,
		message.Cli:      cli.NewCliConnector,
	}
}

type multiChannelConnector struct {
	connector  connector.Connector
	connectors GetConnectors
}

func (m *multiChannelConnector) getConnector(conn message.MessageConnector) (NewConnector, error) {
	newConnector, ok := m.connectors()[conn]
	if !ok {
		return nil, &NotSupportConnectorError{}
	}

	return newConnector, nil
}

func NewMuitiChannelConnector(bot *bot.Bot) *multiChannelConnector {
	m := &multiChannelConnector{}
	m.connectors = Connectors
	connfn, err := m.getConnector(bot.ChatConnector)

	if err != nil {
		logger.Logger.Fatal("chat connector error",
			zap.String("context", "connector"),
			zap.Error(err),
			zap.String("connector_name", string(bot.ChatConnector)))
	}

	conn, err := connfn(bot)

	if err != nil {
		logger.Logger.Fatal("chat connector error",
			zap.String("context", "connector"),
			zap.Error(err),
			zap.String("connector_name", string(bot.ChatConnector)))
	}

	m.connector = conn
	return m
}

func (c *multiChannelConnector) Request(ctx context.Context, input chan<- message.Message) {
	c.connector.Acquire(ctx, input)
}

func (c *multiChannelConnector) Response(ctx context.Context, output <-chan message.Message) {
	for {
		select {
		case msg, ok := <-output:
			if !ok {
				return
			}
			c.connector.Dispatch(msg)
		case <-ctx.Done():
			return
		}
	}
}
