package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	chatCtx "oh-my-chat/src/context"
	"oh-my-chat/src/message"
)

type FakeEgine1 struct {
	engineName string
}

func (f *FakeEgine1) HandleMessage(ctx context.Context, input *message.Message, output chan<- message.Message) {
	input.Output = "message processed"
	output <- *input
}
func (f *FakeEgine1) Config(channelName string) error { return nil }
func (f *FakeEgine1) IsReady() bool                   { return true }

type FakeChatBotGetter struct{}

type Status int

const (
	Processed Status = iota
	ChatNotFound
	EngineNotFound
)

func TestProcess(t *testing.T) {

	assert := assert.New(t)

	type testCase struct {
		desc       string
		status     Status
		engineName string
		botName    string
	}

	for _, scenario := range []testCase{
		{desc: "all ok", status: Processed, engineName: "engine_test", botName: "bot_test"},
		{desc: "chat not found", status: ChatNotFound, engineName: "engine_test", botName: "chat_error"},
		{desc: "engine not found", status: EngineNotFound, engineName: "engine error", botName: "bot_test"},
	} {
		t.Run(scenario.desc, func(t *testing.T) {

			inputMsg := make(chan message.Message, 1)
			outputMsg := make(chan message.Message, 1)
			engine := &FakeEgine1{engineName: scenario.engineName}
			processor := NewProcessor(engine)
			ctx := chatCtx.NewChatContext()
			defer ctx.Shutdown()
			go processor.Process(ctx, inputMsg, outputMsg)
			inputMsg <- message.Message{ID: "123", Input: "hello world", BotName: scenario.botName}

			result := <-outputMsg
			assert.Equal(result.ID, "123")

			switch scenario.status {
			case Processed:
				assert.Equal(result.Output, "message processed")
			case ChatNotFound:
				assert.Equal(result.Output, "some error has ocurred")
				assert.Equal(result.Error, "chat not found")
			case EngineNotFound:
				assert.Equal(result.Output, "some error has ocurred")
				assert.Equal(result.Error, "engine not found")
			}
		})
	}

}
