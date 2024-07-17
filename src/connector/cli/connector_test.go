package cli

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

func TestMain(m *testing.M) {
	logger.InitLog("disable")

	m.Run()
}

type MockBot struct {
	updates int
	mock.Mock
}

func (bot *MockBot) GetUpdateChanels() UpdateChannel {
	update := make(chan Update, 0)

	go func() {
		defer close(update)
		for i := 1; i <= bot.updates; i++ {
			update <- Update{
				UpdateID: i,
				Message:  &Message{Text: "Text", BotName: "test"},
			}
		}
	}()

	return update
}

func (bot *MockBot) SendMessage(message Message) {
	bot.Called(message)
}

func TestAcquire(t *testing.T) {
	assert := assert.New(t)

	t.Run("Ctx done", func(t *testing.T) {
		conn := &cliConnector{&MockBot{updates: 5}}
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			select {
			case <-time.After(time.Millisecond * 100):
				cancel()
			}
		}()

		input := make(chan models.Message, 1)

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-input:
				default:
				}
			}
		}()

		conn.Acquire(ctx, input)
		assert.True(true)
	})

	t.Run("Acquire message", func(t *testing.T) {
		conn := &cliConnector{&MockBot{updates: 5}}
		ctx, cancel := context.WithCancel(context.Background())

		input := make(chan models.Message, 1)
		go func() {
			var receveid int8
			for msg := range input {
				receveid++
				if receveid == 5 {
					cancel()
				}

				assert.Equal("CLI", msg.BotID)
				assert.Equal("Text", msg.Input)
			}
		}()

		conn.Acquire(ctx, input)

	})

}

func TestSendMessage(t *testing.T) {

	t.Run("TextResponse", func(t *testing.T) {
		msg := models.NewMessage()
		msg.ResponseType = models.TextResponse
		msg.Output = "My message"
		mockBot := &MockBot{}
		mockBot.Mock.On("SendMessage", mock.MatchedBy(func(message Message) bool {
			return message.Text == msg.Output && message.IsMultiChoice() == false
		})).Return()

		conn := &cliConnector{bot: mockBot}
		conn.Dispatch(msg)

		mockBot.AssertCalled(t, "SendMessage", mock.AnythingOfType("Message"))

	})

	t.Run("OptionResponse", func(t *testing.T) {
		msg := models.NewMessage()
		msg.ResponseType = models.OptionResponse
		msg.Options = []models.Option{{ID: "1", Name: "Name"}}

		mockBot := &MockBot{}

		mockBot.On("SendMessage", mock.MatchedBy(func(message Message) bool {
			return message.IsMultiChoice() == true && message.MultiChoice[0] == "1"
		})).Return()

		conn := &cliConnector{bot: mockBot}
		conn.Dispatch(msg)

		mockBot.AssertCalled(t, "SendMessage", mock.AnythingOfType("Message"))

	})

}
