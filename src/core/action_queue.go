package core

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

type ActionReplyPair struct {
	replyTo chan<- models.Message
	action  Action
	input   models.Message
}

type goActionQueue struct {
	actionPair chan ActionReplyPair
}

//TODO backpressure strategy

func NewGoActionQueue(workPool int) *goActionQueue {
	return &goActionQueue{actionPair: make(chan ActionReplyPair, workPool)}
}

func (q *goActionQueue) Put(actionPair ActionReplyPair) {
	q.actionPair <- actionPair
}
func (q *goActionQueue) Consume(ctx context.Context) {

	go func() {
		for {
			select {

			case actionPair := <-q.actionPair:
				err := actionPair.action.Handle(ctx, &actionPair.input)

				if err != nil {
					logger.Logger.Error("Error Handling Action",
						zap.String("context", "goActionQueue"),
						zap.Error(err),
					)
				}

				actionPair.replyTo <- actionPair.input

			case <-ctx.Done():
				fmt.Println("Context done")
				q.brodcastAll()

			default:
			}
		}
	}()
}

func (q *goActionQueue) brodcastAll() {
	for {
		select {
		case actionPair := <-q.actionPair:
			actionPair.input.Output = "Server is shutting down. Please reconnect later"

			logger.Logger.Warn("Shutting Down",
				zap.String("context", "goActionQueue"),
				zap.String("message", "context done"),
			)

			actionPair.replyTo <- actionPair.input
		default:
			return
		}
	}
}
