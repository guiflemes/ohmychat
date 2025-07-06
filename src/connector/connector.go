package connector

import (
	"oh-my-chat/src/context"

	"oh-my-chat/src/message"
)

type Connector interface {
	Acquire(ctx *context.ChatContext, input chan<- message.Message)
	Dispatch(message message.Message)
}
