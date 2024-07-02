package core

import (
	"context"

	"go.uber.org/zap"

	"oh-my-chat/src/connector"
	"oh-my-chat/src/connector/cli"
	"oh-my-chat/src/connector/telegram"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

type NotSupportConnectorError struct{}

func (*NotSupportConnectorError) Error() string {
	return "chat connect not supported"
}

type NewConnector func(bot *models.Bot) (connector.Connector, error)
type GetConnectors func() map[models.MessageConnector]NewConnector

func Connectors() map[models.MessageConnector]NewConnector {
	return map[models.MessageConnector]NewConnector{
		models.Telegram: telegram.NewTelegramConnector,
		models.Cli:      cli.NewCliConnector,
	}
}

type multiChannelConnector struct {
	connector  connector.Connector
	connectors GetConnectors
}

func (m *multiChannelConnector) getConnector(conn models.MessageConnector) (NewConnector, error) {
	newConnector, ok := m.connectors()[conn]
	if !ok {
		return nil, &NotSupportConnectorError{}
	}

	return newConnector, nil
}

func NewMuitiChannelConnector(bot *models.Bot) *multiChannelConnector {
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

func (c *multiChannelConnector) Request(ctx context.Context, input chan<- models.Message) {
	c.connector.Acquire(ctx, input)
}

func (c *multiChannelConnector) Response(ctx context.Context, output <-chan models.Message) {
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
