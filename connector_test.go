package ohmychat_test

import (
	"sync"
	"testing"
	"time"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMultiChannelConnector(t *testing.T) {
	t.Parallel()

	t.Run("calls Acquire on connector when Request is invoked", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConnector := mocks.NewMockConnector(ctrl)
		event := make(chan ohmychat.Event)
		chatCtx := ohmychat.NewChatContext(event)
		input := make(chan ohmychat.Message)

		mockConnector.EXPECT().
			Acquire(chatCtx, input).
			Return(nil).
			Times(1)

		mc := ohmychat.NewMuitiChannelConnector(mockConnector)
		mc.Request(chatCtx, input)
	})

	t.Run("calls Acquire on connector when Request is invoked with error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConnector := mocks.NewMockConnector(ctrl)
		event := make(chan ohmychat.Event, 1)
		chatCtx := ohmychat.NewChatContext(event)
		input := make(chan ohmychat.Message)

		mockConnector.EXPECT().
			Acquire(chatCtx, input).
			Return(assert.AnError).
			Times(1)

		mc := ohmychat.NewMuitiChannelConnector(mockConnector)
		go mc.Request(chatCtx, input)

		select {
		case evt := <-event:
			assert.Error(t, evt.Error)
			assert.Nil(t, evt.Msg)
		case <-time.After(200 * time.Millisecond):
			t.Fatal("expected event, but none was received")
		}

		chatCtx.Shutdown()
	})

	t.Run("dispatches messages from output channel", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConnector := mocks.NewMockConnector(ctrl)
		event := make(chan ohmychat.Event, 2)
		chatCtx := ohmychat.NewChatContext(event)
		defer chatCtx.Shutdown()

		output := make(chan ohmychat.Message, 2)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-event
			<-event
			close(event)
		}()

		msg1 := ohmychat.Message{User: ohmychat.User{ID: "a"}, Input: "Hello"}
		msg2 := ohmychat.Message{User: ohmychat.User{ID: "b"}, Input: "World"}

		mockConnector.EXPECT().
			Dispatch(msg1).
			Return(nil).
			Times(1)

		mockConnector.EXPECT().
			Dispatch(msg2).
			Return(nil).
			Times(1)

		mc := ohmychat.NewMuitiChannelConnector(mockConnector)

		output <- msg1
		output <- msg2
		close(output)

		mc.Response(chatCtx, output)

		wg.Wait()
	})

	t.Run("stops dispatching when context is done", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConnector := mocks.NewMockConnector(ctrl)

		event := make(chan ohmychat.Event)
		chatCtx := ohmychat.NewChatContext(event)

		output := make(chan ohmychat.Message)

		mc := ohmychat.NewMuitiChannelConnector(mockConnector)

		chatCtx.Shutdown()

		done := make(chan struct{})

		go func() {
			mc.Response(chatCtx, output)
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Error("Response did not return after context shutdown")
		}
	})
}
