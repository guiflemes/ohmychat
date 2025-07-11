//go:generate mockgen -source processor.go -destination ./mocks/processor.go -package mocks
package core

import (
	"log"
	"oh-my-chat/src/message"
)

type Engine interface {
	HandleMessage(*Context, *message.Message)
}

type processor struct {
	engine Engine
}

func NewProcessor(engine Engine) *processor {
	return &processor{
		engine: engine,
	}
}

func (p *processor) Process(
	ctx *ChatContext,
	inputMsg <-chan message.Message,
	outputMsg chan<- message.Message,
) {
	for {
		select {
		case msg, ok := <-inputMsg:
			if !ok {
				return
			}
			go func(m message.Message) {
				childCtx, err := ctx.NewChildContext(m, outputMsg)
				if err != nil {
					log.Printf("error creating childCtx for session %s", m.User.ID)
					return
				}

				p.engine.HandleMessage(childCtx, &m)

				if !childCtx.MessageHasBeenReplyed() {
					ctx.SaveSession(childCtx.Context(), childCtx.Session())
				}

				childCtx.Cancel()
			}(msg)

		case <-ctx.Done():
			return
		}
	}
}
