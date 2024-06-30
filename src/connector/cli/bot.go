package cli

import (
	"time"

	"github.com/abiosoft/ishell"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

func NewCliBot(botConfig *models.Bot) *CliBot {
	shell := ishell.New()

	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count >= 1 {
			c.Println("Interrupted")
			shell.Stop()
		}
		c.Println("Input Ctrl-c once more to exit")
	})

	go func() {
		shell.Run()
	}()

	cliBot := newCliBot(func() []string { return []string{"cli_test"} })

	shell.AddCmd(&ishell.Cmd{
		Name: "chat",
		Help: "start ohmychat",
		Func: func(c *ishell.Context) {
			cliBot.StartChat(c)
		},
	})

	return cliBot
}

func NewMessage(text string) Message {
	return Message{MessageID: 1, Date: time.Now(), Text: text}
}

type Message struct {
	BotName         string
	MessageID       int
	Date            time.Time
	Text            string
	MultiChoice     []string
	UnBlockByAction bool
}

func (m Message) IsMultiChoice() bool {
	return len(m.MultiChoice) > 0
}

type Update struct {
	UpdateID int
	Message  *Message
}

type UpdateChannel <-chan Update

type ListWorflows func() []string

type CliBot struct {
	Buffer          int
	shutdownChannel chan struct{}
	receiveCh       chan string
	shellCtxt       *ishell.Context
	multiChoiceCh   chan Message
	blocked         bool
	workflow        string
	listWorflows    ListWorflows
}

func newCliBot(listWorflows ListWorflows) *CliBot {

	return &CliBot{
		Buffer:          10,
		shutdownChannel: make(chan struct{}, 1),
		receiveCh:       make(chan string, 10),
		multiChoiceCh:   make(chan Message, 1),
		blocked:         false,
		listWorflows:    listWorflows,
	}
}

func (bot *CliBot) IsRunning() bool {
	return bot.shellCtxt != nil
}

func (bot *CliBot) StartChat(c *ishell.Context) {
	bot.shellCtxt = c

	workflows := bot.listWorflows()
	choice := c.MultiChoice(workflows, "select a workflow")
	bot.workflow = workflows[choice]

	for {
		select {

		case message, ok := <-bot.multiChoiceCh:
			if ok {
				choice := bot.shellCtxt.MultiChoice(message.MultiChoice, "select your choice:")
				bot.receiveCh <- message.MultiChoice[choice]
				break
			}
		default:
			if bot.blocked {
				continue
			}
			bot.blocked = true

			bot.shellCtxt.Print("You: ")
			input := bot.shellCtxt.ReadLine()

			if input == "" {
				continue
			}

			if input == "exit" {
				bot.shellCtxt.Println("Exiting chat mode...")
				return
			}
			bot.receiveCh <- input

		}
	}
}

func (bot *CliBot) GetUpdateChanels() UpdateChannel {
	ch := make(chan Update, bot.Buffer)

	go func() {
		for {
			select {
			case <-bot.shutdownChannel:
				close(ch)
			case receive := <-bot.receiveCh:

				ch <- Update{
					UpdateID: 1,
					Message: &Message{
						BotName:   bot.workflow,
						MessageID: 0,
						Date:      time.Now(),
						Text:      receive,
					},
				}

			}
		}
	}()

	return ch
}

func (bot *CliBot) StopReceivingUpdates() {
	close(bot.shutdownChannel)
}

func (bot *CliBot) SendMessage(message Message) {

	if !bot.IsRunning() {
		logger.Logger.Warn("message dit not send, shell ctx is nil")
		return
	}

	if message.IsMultiChoice() {
		bot.multiChoiceCh <- message
		return
	}

	bot.shellCtxt.Println(message.Text)
	if message.UnBlockByAction {
		bot.blocked = false
	}

}
