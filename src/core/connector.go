package core

import (
	"oh-my-chat/src/message"
)

type Connector interface {
	Acquire(ctx *ChatContext, input chan<- message.Message)
	Dispatch(message message.Message)
}

type multiChannelConnector struct {
	connector Connector
}

func NewMuitiChannelConnector(conn Connector) *multiChannelConnector {
	return &multiChannelConnector{connector: conn}
}

func (c *multiChannelConnector) Request(ctx *ChatContext, input chan<- message.Message) {
	c.connector.Acquire(ctx, input)
}

func (c *multiChannelConnector) Response(ctx *ChatContext, output <-chan message.Message) {
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
