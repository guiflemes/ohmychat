//go:generate mockgen -source processor.go -destination ./mocks/processor.go -package mocks
package core

import (
	"github.com/guiflemes/ohmychat/src/message"
)

type ProcessConfig struct {
	MaxPool uint8
}

type Engine interface {
	HandleMessage(*Context, *message.Message)
}

type processor struct {
	config ProcessConfig
	engine Engine
}

func NewProcessor(engine Engine) *processor {
	return &processor{
		engine: engine,
		config: ProcessConfig{MaxPool: 5},
	}
}

func (p *processor) Process(
	ctx *ChatContext,
	inputMsg <-chan message.Message,
	outputMsg chan<- message.Message,
) {
	sem := make(chan struct{}, p.config.MaxPool)
	for {
		select {
		case msg, ok := <-inputMsg:
			if !ok {
				return
			}
			go func(m message.Message) {
				sem <- struct{}{}
				defer func() { <-sem }()

				childCtx, err := ctx.NewChildContext(m, outputMsg)
				if err != nil {
					ctx.SendEvent(NewEventErrorWithMessage(msg, err))
					return
				}

				p.engine.HandleMessage(childCtx, &m)

				if !childCtx.MessageHasBeenReplyed() {
					if err = ctx.SaveSession(childCtx.Context(), childCtx.Session()); err != nil {
						ctx.SendEvent(NewEventErrorWithMessage(msg, err))
					}
				}

				childCtx.Cancel()
			}(msg)

		case <-ctx.Done():
			return
		}
	}
}
