//go:generate mockgen -source connector.go -destination ./mocks/connector.go -package mocks
package core

import (
	"oh-my-chat/src/message"
)

type Connector interface {
	Acquire(ctx *ChatContext, input chan<- message.Message) error
	Dispatch(message message.Message) error
}

type ConnectorConfig struct {
	ResponseMaxPool uint8
}

type multiChannelConnector struct {
	config    ConnectorConfig
	connector Connector
}

func NewMuitiChannelConnector(conn Connector) *multiChannelConnector {
	return &multiChannelConnector{connector: conn, config: ConnectorConfig{ResponseMaxPool: 5}}
}

func (c *multiChannelConnector) Request(ctx *ChatContext, input chan<- message.Message) {
	err := c.connector.Acquire(ctx, input)
	if err != nil {
		ctx.SendEvent(NewEventError(err))
	}
}

func (c *multiChannelConnector) Response(ctx *ChatContext, output <-chan message.Message) {
	sem := make(chan struct{}, c.config.ResponseMaxPool)
	for {
		select {
		case msg, ok := <-output:
			if !ok {
				return
			}
			go func(m message.Message) {
				sem <- struct{}{}
				defer func() { <-sem }()
				err := c.connector.Dispatch(msg)

				event := NewEvent(msg)
				if err != nil {
					event.WithError(err)
				}
				ctx.SendEvent(*event)

			}(msg)
		case <-ctx.Done():
			return
		}
	}
}
