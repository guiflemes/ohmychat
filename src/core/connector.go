package core

import (
	"go.uber.org/zap"

	"oh-my-chat/src/connector"
	"oh-my-chat/src/connector/telegram"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

type NotSupportConnectorError struct{}

func (*NotSupportConnectorError) Error() string {
	return "chat connect not supported"
}

type NewConnector func(bot *models.Bot) (connector.Connector, error)

type multiChannelConnector struct {
	connector connector.Connector
}

func (m *multiChannelConnector) getConnector(bot *models.Bot) (connector.Connector, error) {
	newConnector, ok := m.Connectors()[bot.ChatConnector]
	if !ok {
		return nil, &NotSupportConnectorError{}
	}

	return newConnector(bot)
}

func NewMuitiChannelConnector(bot *models.Bot) *multiChannelConnector {
	m := &multiChannelConnector{}
	conn, err := m.getConnector(bot)
	if err != nil {
		logger.Logger.Panic("chat connector error",
			zap.String("context", "connector"),
			zap.Error(err),
			zap.String("connector_name", string(bot.ChatConnector)))
		panic("chat connector is not supported")
	}
	m.connector = conn
	return m
}

func (c *multiChannelConnector) Connectors() map[models.MessageConnector]NewConnector {
	return map[models.MessageConnector]NewConnector{
		models.Telegram: telegram.NewTelegramConnector,
	}
}

func (c *multiChannelConnector) Request(input chan<- models.Message) {
	c.connector.Acquire(input)
}

func (c *multiChannelConnector) Respose(output <-chan models.Message) {
	for {
		msg := <-output
		c.connector.Dispatch(msg)
	}
}
