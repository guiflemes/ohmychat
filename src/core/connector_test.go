package core_test

import (
	"testing"
	"time"

	"oh-my-chat/src/core"
	"oh-my-chat/src/core/mocks"
	"oh-my-chat/src/message"

	"github.com/golang/mock/gomock"
)

func TestMultiChannelConnector(t *testing.T) {
	t.Parallel()

	t.Run("calls Acquire on connector when Request is invoked", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConnector := mocks.NewMockConnector(ctrl)
		chatCtx := core.NewChatContext()
		input := make(chan message.Message)

		mockConnector.EXPECT().
			Acquire(chatCtx, input).
			Times(1)

		mc := core.NewMuitiChannelConnector(mockConnector)
		mc.Request(chatCtx, input)
	})

	t.Run("dispatches messages from output channel", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConnector := mocks.NewMockConnector(ctrl)
		chatCtx := core.NewChatContext()
		defer chatCtx.Shutdown()

		output := make(chan message.Message, 2)

		msg1 := message.Message{User: message.User{ID: "a"}, Input: "Hello"}
		msg2 := message.Message{User: message.User{ID: "b"}, Input: "World"}

		mockConnector.EXPECT().
			Dispatch(msg1).
			Times(1)

		mockConnector.EXPECT().
			Dispatch(msg2).
			Times(1)

		mc := core.NewMuitiChannelConnector(mockConnector)

		output <- msg1
		output <- msg2
		close(output)

		mc.Response(chatCtx, output)
	})

	t.Run("stops dispatching when context is done", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockConnector := mocks.NewMockConnector(ctrl)
		chatCtx := core.NewChatContext()

		output := make(chan message.Message)

		mc := core.NewMuitiChannelConnector(mockConnector)

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
