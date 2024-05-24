package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go.uber.org/zap"

	"oh-my-chat/src/core"
	"oh-my-chat/src/logger"
)

type Queue interface {
	Dequeue() ([]byte, bool)
}

type Worker struct {
	queue Queue
}

func (w *Worker) unmarshallMessage(message []byte) core.ActionReplyPair {
	var v core.ActionReplyPair
	if err := json.Unmarshal(message, &v); err != nil {
		v.Input.Output = "some error has ocurred"
		v.Input.Error = ""
		logger.Logger.Error("Error unmarshallMessage",
			zap.String("context", "worker"),
			zap.Error(err))
		return v
	}

	return v

}

func (w *Worker) Produce(ctx context.Context, action chan<- core.ActionReplyPair) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Context done")
			return
		default:
			message, ok := w.queue.Dequeue()
			if ok {
				actionPair := w.unmarshallMessage(message)
				action <- actionPair
			}

		}
	}
}

func (w *Worker) Consume(ctx context.Context, action <-chan core.ActionReplyPair) {

	for {
		select {

		case <-ctx.Done():
			fmt.Println("Context done")
			return

		case actionPair := <-action:
			err := actionPair.Action.Handle(ctx, &actionPair.Input)
			if err != nil {
				logger.Logger.Error("Error Handling Action",
					zap.String("context", "Worker"),
					zap.Error(err),
				)
			}

			actionPair.ReplyTo <- actionPair.Input

		default:
		}
	}
}

func RunWorker(ctx context.Context, queue Queue) {
	actionCh := make(chan core.ActionReplyPair)
	var wg sync.WaitGroup

	worker := &Worker{queue: queue}

	wg.Add(2)
	go worker.Produce(ctx, actionCh)
	go worker.Consume(ctx, actionCh)

	wg.Wait()
}
