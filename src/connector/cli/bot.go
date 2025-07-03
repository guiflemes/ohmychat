package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/abiosoft/ishell"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

func NewCliBot(botConfig *models.Bot, shell *ishell.Shell) *CliBot {

	listWorflows := botConfig.CliDependencies.ListWorkflows

	if listWorflows == nil {
		listWorflows = func() []string { return []string{"new_chat"} }
	}

	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count >= 1 {
			c.Println("Interrupted")
			shell.Stop()
		}
		c.Println("Input Ctrl-c once more to exit")
	})

	go func() {
		if !botConfig.CliDependencies.DisableInitialization {

			fmt.Println(`
     ( )
.-----'-----.
| ( )   ( ) |  -( welcome to ohmychat !!! type 'chat' to start it or 'help' to see all options )
'-----.-----' 
 / '+---+' \ 
 \/--|_|--\/`)
		}
		shell.Run()
	}()

	cliBot := newCliBot(listWorflows())

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

type ListWorflows []string

type CliBot struct {
	Buffer          int
	shutdownChannel chan struct{}
	receiveCh       chan string
	shellCtx        *ishell.Context
	multiChoiceCh   chan Message
	outputCh        chan Message
	blocked         bool
	workflow        string
	listWorflows    ListWorflows
	waitingResponse bool
}

func newCliBot(listWorflows ListWorflows) *CliBot {

	return &CliBot{
		Buffer:          10,
		shutdownChannel: make(chan struct{}, 1),
		receiveCh:       make(chan string, 10),
		multiChoiceCh:   make(chan Message, 1),
		outputCh:        make(chan Message, 1),
		blocked:         false,
		listWorflows:    listWorflows,
		waitingResponse: false,
	}
}

func (bot *CliBot) IsRunning() bool {
	return bot.shellCtx != nil
}

func (bot *CliBot) StartChat(c *ishell.Context) {
	bot.shellCtx = c

	choice := c.MultiChoice(bot.listWorflows, "select a workflow")
	bot.workflow = bot.listWorflows[choice]

	for {
		if !bot.waitingResponse {
			bot.shellCtx.Print("YOU: ")
			input := bot.shellCtx.ReadLine()
			input = strings.TrimSpace(input)

			if input == "exit" {
				bot.shellCtx.Println("Exiting chat mode...")
				return
			}

			bot.receiveCh <- input
			bot.waitingResponse = true
		}

		select {
		case message := <-bot.multiChoiceCh:
			choice := bot.shellCtx.MultiChoice(message.MultiChoice, "select your choice:")
			bot.receiveCh <- message.MultiChoice[choice]
		case message := <-bot.outputCh:
			bot.shellCtx.Print("BOT: ")
			bot.shellCtx.Println(message.Text)
			bot.waitingResponse = false
			// default:
			// 	// if bot.blocked {
			// 	// 	time.Sleep(100 * time.Millisecond)
			// 	// 	continue
			// 	// }
			// 	if bot.waitingResponse {
			// 		continue
			// 	}

			// 	bot.shellCtx.Print("YOU: ")
			// 	input := bot.shellCtx.ReadLine()
			// 	input = strings.TrimSpace(input)

			// 	if input == "exit" {
			// 		bot.shellCtx.Println("Exiting chat mode...")
			// 		return
			// 	}
			// 	bot.receiveCh <- input

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
	bot.shutdownChannel <- struct{}{}
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

	bot.outputCh <- message

	if message.UnBlockByAction {
		bot.blocked = false
	}

}
