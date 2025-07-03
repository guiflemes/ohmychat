package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"oh-my-chat/src/models"
)

type FakeEgine1 struct {
	engineName string
}

func (f *FakeEgine1) Name() string { return f.engineName }
func (f *FakeEgine1) HandleMessage(ctx context.Context, input models.Message, output chan<- models.Message) {
	response := &input
	response.Output = "message processed"
	output <- *response
}
func (f *FakeEgine1) GetActionStorageService() ActionStorageService { return nil }
func (f *FakeEgine1) Config(channelName string) error               { return nil }
func (f *FakeEgine1) IsReady() bool                                 { return true }

type FakeChatBotGetter struct{}

func (f *FakeChatBotGetter) GetChatBot(botName string) *models.ChatBot {
	if botName != "bot_test" {
		return nil
	}
	return &models.ChatBot{
		BotName:    "bot_test",
		Engine:     "engine_test",
		WorkflowID: "workflow_test",
	}
}

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

			inputMsg := make(chan models.Message, 1)
			outputMsg := make(chan models.Message, 1)
			chatBotRepo := &FakeChatBotGetter{}
			engine := &FakeEgine1{engineName: scenario.engineName}
			processor := NewProcessor(chatBotRepo, []Engine{engine})

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go processor.Process(ctx, inputMsg, outputMsg)
			inputMsg <- models.Message{ID: "123", Input: "hello world", BotName: scenario.botName}

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
