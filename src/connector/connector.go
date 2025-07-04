package connector

import (
	"context"

	"oh-my-chat/src/message"
)

type Connector interface {
	Acquire(ctx context.Context, input chan<- message.Message)
	Dispatch(message message.Message)
}
