package ohmychat_test

import (
	"testing"
	"time"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChatContextAndContext(t *testing.T) {
	t.Parallel()

	t.Run("creates context with default session adapter", func(t *testing.T) {
		t.Parallel()

		ctx := ohmychat.NewChatContext(make(chan<- ohmychat.Event))
		assert.True(t, ctx.IsActive())

		ctx.Set("foo", "bar")
		val, ok := ctx.Get("foo")
		assert.True(t, ok)
		assert.Equal(t, "bar", val)

		ctx.Shutdown()
		assert.False(t, ctx.IsActive())
	})

	t.Run("creates context with custom session adapter and returns child context", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSessionAdapter(ctrl)
		session := &ohmychat.Session{UserID: "abc", Memory: make(map[string]any), State: ohmychat.IdleState{}}

		mockAdapter.EXPECT().
			GetOrCreate(gomock.Any(), "abc").
			Return(session, nil).
			Times(1)

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdapter),
		)

		msg := ohmychat.Message{User: ohmychat.User{ID: "abc"}}
		output := make(chan ohmychat.Message, 1)

		child, err := chatCtx.NewChildContext(msg, output)
		assert.NoError(t, err)
		assert.NotNil(t, child)
		assert.Equal(t, session, child.Session())
	})

	t.Run("send output sends message and saves session", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockAdapter := mocks.NewMockSessionAdapter(ctrl)
		session := &ohmychat.Session{UserID: "kizaru", Memory: make(map[string]any), State: ohmychat.IdleState{}}

		mockAdapter.EXPECT().
			GetOrCreate(gomock.Any(), "kizaru").
			Return(session, nil).
			Times(1)

		mockAdapter.EXPECT().
			Save(gomock.Any(), session).
			Return(nil).
			Times(1)

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdapter),
		)

		msg := ohmychat.Message{User: ohmychat.User{ID: "kizaru"}}
		output := make(chan ohmychat.Message, 1)

		childCtx, err := chatCtx.NewChildContext(msg, output)
		assert.NoError(t, err)

		toSend := &ohmychat.Message{User: ohmychat.User{ID: "kizaru"}, Input: "hello!"}
		childCtx.SendOutput(toSend)

		select {
		case received := <-output:
			assert.Equal(t, toSend.Input, received.Input)
			assert.Equal(t, toSend.User.ID, received.User.ID)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("expected message on output channel")
		}
	})
}
