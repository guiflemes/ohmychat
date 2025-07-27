package ohmychat_test

import (
	"testing"
	"time"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestProcessor_Process(t *testing.T) {
	t.Parallel()

	t.Run("sucessfully", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEngine := mocks.NewMockEngine(ctrl)

		mockEngine.EXPECT().HandleMessage(gomock.Any(), gomock.Any()).Times(1)

		proc := ohmychat.NewProcessor(mockEngine)

		input := make(chan ohmychat.Message, 1)
		output := make(chan ohmychat.Message, 1)
		event := make(chan ohmychat.Event, 1)

		msg := ohmychat.Message{
			User: ohmychat.User{ID: "user123"},
		}

		ctx := ohmychat.NewChatContext(event)

		input <- msg

		go proc.Process(ctx, input, output)

		time.Sleep(100 * time.Millisecond)

		ctx.Shutdown()
		close(event)

	})

	t.Run("skip handling message, session adapter error", func(t *testing.T) {
		t.Parallel()
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockEngine := mocks.NewMockEngine(ctrl)
		mockSessionAdapter := mocks.NewMockSessionAdapter(ctrl)

		mockEngine.EXPECT().HandleMessage(gomock.Any(), gomock.Any()).Times(0)
		mockSessionAdapter.EXPECT().GetOrCreate(gomock.Any(), "user123").Return(nil, assert.AnError).Times(1)

		proc := ohmychat.NewProcessor(mockEngine)

		input := make(chan ohmychat.Message, 1)
		output := make(chan ohmychat.Message, 1)
		event := make(chan ohmychat.Event, 1)

		msg := ohmychat.Message{
			User: ohmychat.User{ID: "user123"},
		}

		ctx := ohmychat.NewChatContext(
			event,
			ohmychat.WithSessionAdapter(mockSessionAdapter),
		)

		input <- msg

		go proc.Process(ctx, input, output)

		select {
		case evt := <-event:
			assert.Error(t, evt.Error)
			assert.Equal(t, "user123", evt.Msg.User.ID)
			assert.ErrorIs(t, assert.AnError, evt.Error)
		case <-time.After(200 * time.Millisecond):
			t.Fatal("expected event, but none was received")
		}
		ctx.Shutdown()

	})
}
