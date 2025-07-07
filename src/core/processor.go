package core

import (
	"context"
	"oh-my-chat/src/message"

	chatCtx "oh-my-chat/src/context"
)

type Engine interface {
	HandleMessage(context.Context, *message.Message, chan<- message.Message)
}

type processor struct {
	engine Engine
}

func NewProcessor(engine Engine) *processor {
	return &processor{
		engine: engine,
	}
}

func (m *processor) Process(
	ctx *chatCtx.ChatContext,
	inputMsg <-chan message.Message,
	outputMsg chan<- message.Message,
) {
	for {
		select {
		case message, ok := <-inputMsg:
			if !ok {
				return
			}
			childCtx := ctx.NewChildContext()
			m.engine.HandleMessage(childCtx.Context(), &message, outputMsg)

		case <-ctx.Done():
			return
		}

	}
}
