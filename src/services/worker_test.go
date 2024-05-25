package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"oh-my-chat/src/core"
	"oh-my-chat/src/models"
)

type MockStorage struct {
	data []*core.ActionReplyPair
}

func (m *MockStorage) Pop() (*core.ActionReplyPair, bool) {
	if len(m.data) == 0 {
		return nil, false
	}

	item := m.data[0]
	m.data = m.data[1:]
	return item, true
}

type WorkerSuite struct {
	suite.Suite
	worker *Worker
}

func (w *WorkerSuite) SetupTest() {
	w.worker = &Worker{storage: &MockStorage{
		data: []*core.ActionReplyPair{
			{
				ReplyTo: make(chan<- models.Message),
				Input:   models.NewMessage(),
				Action:  nil,
			},
			{
				ReplyTo: make(chan<- models.Message),
				Input:   models.NewMessage(),
				Action:  nil,
			},
		},
	}}
}

func (w *WorkerSuite) TestProducer() {
	ctx := context.Background()
	ch := make(chan core.ActionReplyPair)

	go func() {
		w.worker.Produce(ctx, ch)
	}()

	messsage1, ok := <-ch

	w.True(ok)
	w.NotNil(messsage1)

	messsage2, ok := <-ch
	w.True(ok)
	w.NotNil(messsage2)

}
func (w *WorkerSuite) TestProducerCtxDone() {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan core.ActionReplyPair)

	go func() {
		w.worker.Produce(ctx, ch)
	}()

	cancel()

	select {
	case msg := <-ch:
		w.Require().Nil(msg, "Expected no messages after context cancel")
	case <-time.After(time.Millisecond * 100):
	}

}

func TestWorkSuite(t *testing.T) {
	suite.Run(t, new(WorkerSuite))
}
