package cli

import (
	"time"

	"github.com/abiosoft/ishell"

	"oh-my-chat/src/logger"
)

func NewMessage(text string) Message {
	return Message{MessageID: 1, Date: time.Now(), Text: text}
}

type Message struct {
	MessageID   int
	Date        time.Time
	Text        string
	MultiChoice []string
}

func (m Message) IsMultiChoice() bool {
	return len(m.MultiChoice) > 0
}

type Update struct {
	UpdateID int
	Message  *Message
}

type UpdateChannel <-chan Update

type CliBot struct {
	Buffer          int
	shutdownChannel chan struct{}
	receiveCh       chan string
	shellCtxt       *ishell.Context
	lastReplyMsg    Message
	multiChoiceCh   chan Message
	blocked         bool
}

func NewCliBot() *CliBot {
	return &CliBot{
		Buffer:          10,
		shutdownChannel: make(chan struct{}, 1),
		receiveCh:       make(chan string, 10),
		multiChoiceCh:   make(chan Message, 1),
		blocked:         false,
	}
}

func (bot *CliBot) IsRunning() bool {
	return bot.shellCtxt != nil
}

func (bot *CliBot) StartChat(c *ishell.Context) {
	bot.shellCtxt = c

	for {
		select {
		case message, ok := <-bot.multiChoiceCh:
			if ok {
				choice := bot.shellCtxt.MultiChoice(message.MultiChoice, "Escolha sua opção:")
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
	bot.lastReplyMsg = message

	if !bot.IsRunning() {
		logger.Logger.Warn("message dit not send, shell ctx is nil")
		return
	}

	if message.IsMultiChoice() {
		bot.multiChoiceCh <- message
		return
	}

	bot.shellCtxt.Println("reply", message.Text)
	bot.blocked = false
}
