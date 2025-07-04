package core

import (
	"context"

	"oh-my-chat/src/models"
)

type Engine interface {
	HandleMessage(context.Context, *models.Message, chan<- models.Message)
}

type processor struct {
	chatBot *models.ChatBot
	engine  Engine
}

func NewProcessor(engine Engine) *processor {
	return &processor{
		engine: engine,
	}
}

func (m *processor) Process(
	ctx context.Context,
	inputMsg <-chan models.Message,
	outputMsg chan<- models.Message,
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
