package connector

import (
	"context"

	"oh-my-chat/src/models"
)

type Connector interface {
	Acquire(ctx context.Context, input chan<- models.Message)
	Dispatch(message models.Message)
}
