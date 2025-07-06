package core

import (
	"context"

	"oh-my-chat/src/connector"
	"oh-my-chat/src/message"
)

type multiChannelConnector struct {
	connector connector.Connector
}

func NewMuitiChannelConnector(conn connector.Connector) *multiChannelConnector {
	return &multiChannelConnector{connector: conn}
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
