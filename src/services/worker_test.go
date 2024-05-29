package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"oh-my-chat/src/config"
	"oh-my-chat/src/core"
	"oh-my-chat/src/models"
)

type MockStorage struct {
	data []core.ActionReplyPair
}

func (m *MockStorage) Dequeue() (core.ActionReplyPair, bool) {
	if len(m.data) == 0 {
		return core.ActionReplyPair{}, false
	}

	item := m.data[0]
	m.data = m.data[1:]
	return item, true
}

type MockAction struct {
	mock.Mock
	called bool
}

func (m *MockAction) Handle(ctx context.Context, message *models.Message) error {
	args := m.Called(ctx, message)
	message.Output = "Output"
	return args.Error(0)
}

type WorkerSuite struct {
	suite.Suite
	worker      *Worker
	actions     []core.ActionReplyPair
	mockActions []*MockAction
	replyToCh   chan models.Message
}

func (w *WorkerSuite) BeforeTest(suiteName, testName string) {
	mockAction1 := &MockAction{}
	mockAction2 := &MockAction{}
	w.mockActions = []*MockAction{mockAction1, mockAction2}
	w.replyToCh = make(chan models.Message)

	w.actions = []core.ActionReplyPair{
		{
			ReplyTo: w.replyToCh,
			Input:   models.NewMessage(),
			Action:  mockAction1,
		},
		{
			ReplyTo: w.replyToCh,
			Input:   models.NewMessage(),
			Action:  mockAction2,
		},
	}

	w.worker = &Worker{storage: &MockStorage{
		data: w.actions,
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
	close(ch)

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

	close(ch)
}
func (w *WorkerSuite) TestConsumer() {
	actionCh := make(chan core.ActionReplyPair)

	w.Run("Context Done", func() {
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			time.Sleep(time.Millisecond * 100)
			defer cancel()
		}()

		w.worker.Consume(ctx, actionCh)
		w.True(true)
	})

	w.Run("Handle Action Ok", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		actionReplyPair := w.actions[0]
		mockAction1 := w.mockActions[0]
		mockAction1.On("Handle", mock.Anything, &actionReplyPair.Input).Return(nil).Once()
		go w.worker.Consume(ctx, actionCh)
		action, _ := w.worker.storage.Dequeue()
		actionCh <- action

		receive := <-w.replyToCh
		mockAction1.AssertExpectations(w.T())
		w.Equal(receive.Output, "Output")
	})

	w.Run("Handle Action Error", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		actionReplyPair := w.actions[1]
		mockAction2 := w.mockActions[1]
		mockAction2.On("Handle", mock.Anything, &actionReplyPair.Input).
			Return(fmt.Errorf("some error has ocurred")).
			Once()
		go w.worker.Consume(ctx, actionCh)
		action, _ := w.worker.storage.Dequeue()
		actionCh <- action

		receive := <-w.replyToCh
		mockAction2.AssertExpectations(w.T())
		w.Equal(receive.Error, "some error has ocurred")
	})

}

func TestWorkSuite(t *testing.T) {
	suite.Run(t, new(WorkerSuite))
}

func TestRunWorker(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	config := config.Worker{Number: 2}
	storage := &MockStorage{}

	go func() {
		time.Sleep(time.Millisecond * 100)
		cancel()
	}()

	RunWorker(ctx, config, storage)
	assert := assert.New(t)
	assert.True(true, "stop worker after cancel ctx")

}
