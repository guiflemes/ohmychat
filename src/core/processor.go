package core

import (
	"context"

	"oh-my-chat/src/bot"
	"oh-my-chat/src/message"
)

type Engine interface {
	HandleMessage(context.Context, *message.Message, chan<- message.Message)
}

type processor struct {
	chatBot *bot.ChatBot
	engine  Engine
}

func NewProcessor(engine Engine) *processor {
	return &processor{
		engine: engine,
	}
}

func (m *processor) Process(
	ctx context.Context,
	inputMsg <-chan message.Message,
	outputMsg chan<- message.Message,
) {
	for {
		select {
		case message, ok := <-inputMsg:
			if !ok {
				return
			}
			m.engine.HandleMessage(ctx, &message, outputMsg)

		case <-ctx.Done():
			return
		}

	}
}
