//go:generate mockgen -source connector.go -destination ./mocks/connector.go -package mocks
package ohmychat

type Connector interface {
	Acquire(ctx *ChatContext, input chan<- Message) error
	Dispatch(message Message) error
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

func (c *multiChannelConnector) Request(ctx *ChatContext, input chan<- Message) {
	err := c.connector.Acquire(ctx, input)
	if err != nil {
		ctx.SendEvent(NewEventError(err))
	}
}

func (c *multiChannelConnector) Response(ctx *ChatContext, output <-chan Message, input chan<- Message) {
	sem := make(chan struct{}, c.config.ResponseMaxPool)
	for {
		select {
		case msg, ok := <-output:
			if !ok {
				return
			}
			go func(m Message) {
				sem <- struct{}{}
				defer func() { <-sem }()
				err := c.connector.Dispatch(msg)
				event := NewEvent(msg)
				if err != nil {
					event.WithError(err)
				}
				ctx.SendEvent(*event)

				if msg.BotMode {
					input <- msg.NewFrom()
				}

			}(msg)
		case <-ctx.Done():
			return
		}
	}
}
