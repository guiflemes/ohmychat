package rule_engine_test

import (
	"testing"
	"time"

	"github.com/guiflemes/ohmychat/core"
	"github.com/guiflemes/ohmychat/core/mocks"
	"github.com/guiflemes/ohmychat/core/rule_engine"
	"github.com/guiflemes/ohmychat/message"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRuleEngine(t *testing.T) {
	t.Parallel()

	t.Run("handle idle state with no match", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ss := &core.Session{State: core.IdleState{}, LastActivityAt: time.Now()}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		msg := message.Message{Input: "rocksdxebec"}
		output := make(chan message.Message, 1)

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

		ss := &core.Session{State: core.IdleState{}, LastActivityAt: time.Now()}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		msg := &message.Message{Input: "hello"}

		rule := rule_engine.Rule{
			Prompts:   []string{"hello"},
			NextState: core.IdleState{},
			Action: func(ctx *core.Context, m *message.Message) {
				called = true
			},
		}

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		output := make(chan message.Message, 1)

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

		ss := &core.Session{State: core.WaitingInputState{PromptEmptyMessage: "akainu"}, LastActivityAt: time.Now()}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		msg := &message.Message{Input: " "}

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		output := make(chan message.Message, 1)

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

		ss := &core.Session{
			LastActivityAt: time.Now(),
			State: core.WaitingInputState{
				Action: func(ctx *core.Context, m *message.Message) {
					called = true
				},
			},
		}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)

		msg := &message.Message{Input: "zoro"}

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		output := make(chan message.Message, 1)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)

		engine := rule_engine.NewRuleEngine()
		engine.HandleMessage(childCtx, msg)

		assert.True(t, called)
	})

	t.Run("handle waiting choice with invalid option", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ss := &core.Session{
			LastActivityAt: time.Now(),
			State: core.WaitingChoiceState{
				Choices:             map[string]core.ActionFunc{},
				Prompt:              "choose an option",
				PromptInvalidOption: "you're wrong",
			},
		}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)
		msg := &message.Message{Input: "invalid_option"}

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		output := make(chan message.Message, 1)
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

		ss := &core.Session{
			LastActivityAt: time.Now(),
			State: core.WaitingChoiceState{
				Choices: map[string]core.ActionFunc{
					"1": func(ctx *core.Context, m *message.Message) {
						called = true
					},
				},
			},
		}

		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		msg := &message.Message{Input: "1"}

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		output := make(chan message.Message, 1)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)
		engine := rule_engine.NewRuleEngine()

		engine.HandleMessage(childCtx, msg)

		assert.True(t, called)
	})

	t.Run("handle unknown state", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ss := &core.Session{
			LastActivityAt: time.Now(),
			State:          nil,
		}
		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		msg := &message.Message{Input: "teste"}

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		output := make(chan message.Message, 1)
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
		ss := &core.Session{State: core.WaitingInputState{}, LastActivityAt: expiredTime}

		mockAdpater := mocks.NewMockSessionAdapter(ctrl)
		mockAdpater.EXPECT().GetOrCreate(gomock.Any(), gomock.Any()).Return(ss, nil)
		mockAdpater.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

		msg := &message.Message{Input: "Monkey D Garp"}
		output := make(chan message.Message, 1)

		chatCtx := core.NewChatContext(
			make(chan<- core.Event),
			core.WithSessionAdapter(mockAdpater),
		)
		childCtx, _ := chatCtx.NewChildContext(*msg, output)

		engine := rule_engine.NewRuleEngine()
		engine.HandleMessage(childCtx, msg)

		_, isIdle := ss.State.(core.IdleState)
		assert.True(t, isIdle, "Session should reset to IdleState after expiration")
	})
}
