package cli

import (
	"testing"

	"github.com/abiosoft/ishell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"oh-my-chat/src/models"
)

func TestNewCliBot(t *testing.T) {

	t.Run("panic on workflow are empty", func(t *testing.T) {
		assert.Panics(t, func() {
			NewCliBot(&models.Bot{CliDependencies: models.CliDependencies{
				ListWorkflows: func() []string { return []string{} },
			}}, ishell.New())
		})
	})

	t.Run("new bot", func(t *testing.T) {
		shell := ishell.New()
		newBot := NewCliBot(&models.Bot{CliDependencies: models.CliDependencies{
			DisableInitialization: true,
			ListWorkflows:         func() []string { return []string{"choice"} },
		}}, shell)

		assert.NotNil(t, newBot)
		assert.Equal(t, len(newBot.listWorflows), 1)
		shell.Stop()

	})

}
func TestGetUpdateChanels(t *testing.T) {
	assert := assert.New(t)

	cliBot := newCliBot(ListWorflows{})

	go func() {
		for i := 0; i > 5; i++ {
			cliBot.receiveCh <- "my message"
		}
		cliBot.StopReceivingUpdates()
	}()

	for update := range cliBot.GetUpdateChanels() {
		assert.Equal("my message", update.Message.Text)

	}

}

func TestBotSendMessage(t *testing.T) {

	assert := assert.New(t)

	t.Run("Bot is not running", func(t *testing.T) {
		cliBot := newCliBot(ListWorflows{})
		cliBot.SendMessage(Message{})

		assert.False(cliBot.IsRunning())

	})

	t.Run("MultiChoice message", func(t *testing.T) {

		cliBot := newCliBot(ListWorflows{})

		go func() {
			cliBot.shellCtx = &ishell.Context{}
			cliBot.SendMessage(Message{MultiChoice: []string{"choice"}})
		}()

		choice := <-cliBot.multiChoiceCh
		assert.True(choice.IsMultiChoice())
		assert.Equal(len(choice.MultiChoice), 1)
		assert.Equal(choice.MultiChoice[0], "choice")

	})

	t.Run("Text message", func(t *testing.T) {
		cliBot := newCliBot(ListWorflows{})
		mockAction := &MockActions{}
		cliBot.blocked = true
		mockAction.Mock.On("Println", "my message").Return()
		cliBot.shellCtx = &ishell.Context{Actions: mockAction}
		cliBot.SendMessage(Message{Text: "my message", UnBlockByAction: true})
		mockAction.AssertCalled(t, "Println", "my message")
		assert.False(cliBot.blocked)

	})

}

func TestStartChat(t *testing.T) {
	assert := assert.New(t)

	t.Run("readline", func(t *testing.T) {
		cliBot := newCliBot(ListWorflows{"choice"})
		mockAction := &MockActions{readLine: "readline"}

		go func() {
			cliBot.StartChat(&ishell.Context{Actions: mockAction})
		}()

		received := <-cliBot.receiveCh
		assert.Equal(received, "readline")

	})

	t.Run("exit", func(t *testing.T) {
		cliBot := newCliBot(ListWorflows{"choice"})
		mockAction := &MockActions{readLine: "exit"}
		mockAction.Mock.On("Println", "Exiting chat mode...").Return()
		cliBot.StartChat(&ishell.Context{Actions: mockAction})
		mockAction.AssertCalled(t, "Println", "Exiting chat mode...")

	})

	t.Run("MultiChoice", func(t *testing.T) {

		cliBot := newCliBot(ListWorflows{"choice"})
		cliBot.blocked = true
		mockAction := &MockActions{readLine: "readline"}

		go func() {
			cliBot.StartChat(&ishell.Context{Actions: mockAction})
		}()

		go func() {
			cliBot.multiChoiceCh <- Message{MultiChoice: []string{"choice"}}
		}()

		received := <-cliBot.receiveCh
		assert.Equal(received, "choice")
	})

}

type MockActions struct {
	mock.Mock
	readLine string
}

func (a *MockActions) ReadLine() string                              { return a.readLine }
func (a *MockActions) ReadLineErr() (string, error)                  { return "", nil }
func (a *MockActions) ReadPassword() string                          { return "" }
func (a *MockActions) ReadPasswordErr() (string, error)              { return "", nil }
func (a *MockActions) ReadMultiLinesFunc(f func(string) bool) string { return "" }
func (a *MockActions) ReadMultiLines(terminator string) string       { return "" }
func (a *MockActions) Println(val ...interface{}) {
	a.Called(val...)
}
func (a *MockActions) Print(val ...interface{})                                  {}
func (a *MockActions) Printf(format string, val ...interface{})                  {}
func (a *MockActions) ShowPaged(text string) error                               { return nil }
func (a *MockActions) MultiChoice(options []string, text string) int             { return 0 }
func (a *MockActions) Checklist(options []string, text string, init []int) []int { return nil }
func (a *MockActions) SetPrompt(prompt string)                                   {}
func (a *MockActions) SetMultiPrompt(prompt string)                              {}
func (a *MockActions) ShowPrompt(show bool)                                      {}
func (a *MockActions) Cmds() []*ishell.Cmd                                       { return nil }
func (a *MockActions) HelpText() string                                          { return "" }
func (a *MockActions) ClearScreen() error                                        { return nil }
func (a *MockActions) Stop()                                                     {}
