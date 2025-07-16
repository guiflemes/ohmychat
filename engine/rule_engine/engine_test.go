package rule_engine_test

import (
	"testing"
	"time"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/engine/rule_engine"
	"github.com/guiflemes/ohmychat/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRuleEngine(t *testing.T) {
	t.Parallel()

	t.Run("handle idle state with no match", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ss := &ohmychat.Session{State: ohmychat.IdleState{}, LastActivityAt: time.Now()}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		msg := ohmychat.Message{Input: "rocksdxebec"}
		output := make(chan ohmychat.Message, 1)

		childCtx, _ := chatCtx.NewChildContext(msg, output)
		engine := rule_engine.NewRuleEngine()
		engine.HandleMessage(childCtx, &msg)

		assert.Equal(t, "desculpe não entendi", msg.Output)
	})

	t.Run("handle idle state with rule match", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		called := false

		ss := &ohmychat.Session{State: ohmychat.IdleState{}, LastActivityAt: time.Now()}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		msg := &ohmychat.Message{Input: "hello"}

		rule := rule_engine.Rule{
			Prompts:   []string{"hello"},
			NextState: ohmychat.IdleState{},
			Action: func(ctx *ohmychat.Context, m *ohmychat.Message) {
				called = true
			},
		}

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		output := make(chan ohmychat.Message, 1)

		childCtx, _ := chatCtx.NewChildContext(*msg, output)

		engine := rule_engine.NewRuleEngine()
		engine.RegisterRule(rule)

		engine.HandleMessage(childCtx, msg)
		assert.True(t, called)
	})

	t.Run("handle waiting input with empty input", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ss := &ohmychat.Session{State: ohmychat.WaitingInputState{PromptEmptyMessage: "akainu"}, LastActivityAt: time.Now()}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		msg := &ohmychat.Message{Input: " "}

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		output := make(chan ohmychat.Message, 1)

		childCtx, _ := chatCtx.NewChildContext(*msg, output)

		engine := rule_engine.NewRuleEngine()

		engine.HandleMessage(childCtx, msg)

		assert.Equal(t, "akainu", msg.Output)
	})

	t.Run("handle waiting input with valid input", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		called := false

		ss := &ohmychat.Session{
			LastActivityAt: time.Now(),
			State: ohmychat.WaitingInputState{
				Action: func(ctx *ohmychat.Context, m *ohmychat.Message) {
					called = true
				},
			},
		}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)

		msg := &ohmychat.Message{Input: "zoro"}

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		output := make(chan ohmychat.Message, 1)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)

		engine := rule_engine.NewRuleEngine()
		engine.HandleMessage(childCtx, msg)

		assert.True(t, called)
	})

	t.Run("handle waiting choice with invalid option", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ss := &ohmychat.Session{
			LastActivityAt: time.Now(),
			State: ohmychat.WaitingChoiceState{
				Choices:             map[string]ohmychat.ActionFunc{},
				Prompt:              "choose an option",
				PromptInvalidOption: "you're wrong",
			},
		}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		msg := &ohmychat.Message{Input: "invalid_option"}

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		output := make(chan ohmychat.Message, 1)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)
		engine := rule_engine.NewRuleEngine()
		engine.HandleMessage(childCtx, msg)

		assert.Equal(t, msg.Output, "you're wrong")
	})

	t.Run("handle waiting choice with valid option", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		called := false

		ss := &ohmychat.Session{
			LastActivityAt: time.Now(),
			State: ohmychat.WaitingChoiceState{
				Choices: map[string]ohmychat.ActionFunc{
					"1": func(ctx *ohmychat.Context, m *ohmychat.Message) {
						called = true
					},
				},
			},
		}

		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		msg := &ohmychat.Message{Input: "1"}

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		output := make(chan ohmychat.Message, 1)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)
		engine := rule_engine.NewRuleEngine()

		engine.HandleMessage(childCtx, msg)

		assert.True(t, called)
	})

	t.Run("handle unknown state", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ss := &ohmychat.Session{
			LastActivityAt: time.Now(),
			State:          nil,
		}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		msg := &ohmychat.Message{Input: "teste"}

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		output := make(chan ohmychat.Message, 1)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)
		engine := rule_engine.NewRuleEngine()

		engine.HandleMessage(childCtx, msg)

		assert.Equal(t, "Erro interno: estado desconhecido.", msg.Output)
	})

	t.Run("handle expired session and reset to idle state", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expiredTime := time.Now().Add(-10 * time.Minute)
		ss := &ohmychat.Session{State: ohmychat.WaitingInputState{}, LastActivityAt: expiredTime}

		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		msg := &ohmychat.Message{Input: "Monkey D Garp"}
		output := make(chan ohmychat.Message, 1)

		chatCtx := ohmychat.NewChatContext(
			make(chan<- ohmychat.Event),
			ohmychat.WithSessionAdapter(mockAdpater),
		)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)

		engine := rule_engine.NewRuleEngine()
		engine.HandleMessage(childCtx, msg)

		_, isIdle := ss.State.(ohmychat.IdleState)
		assert.True(t, isIdle, "Session should reset to IdleState after expiration")
	})
}
