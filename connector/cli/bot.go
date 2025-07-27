package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
)

type CliOption func(cli *CliBot)

func NewCliBot(shell *ishell.Shell, control *ChatControl, opts ...CliOption) *CliBot {

	cliBot := &CliBot{
		Buffer:          10,
		shutdownChannel: make(chan struct{}, 1),
		receiveCh:       make(chan string, 10),
		multiChoiceCh:   make(chan Message, 1),
		outputCh:        make(chan Message, 1),
		waitingResponse: false,
	}

	for _, opt := range opts {
		opt(cliBot)
	}

	if cliBot.listWorflows == nil {
		cliBot.listWorflows = []string{"new_chat"}
	}

	shell.Interrupt(func(c *ishell.Context, count int, input string) {
		if count >= 1 {
			c.Println("Interrupted")
			shell.Stop()
		}
	})

	go func() {
		if !cliBot.disableInitialization {

			fmt.Println(`
     ( )
.-----'-----.
| ( )   ( ) |  -( welcome to ohmychat !!! type 'chat' to start it or 'help' to see all options )
'-----.-----' 
 / '+---+' \ 
 \/--|_|--\/`)
		}
		shell.Run()
		if control.ctx == nil {
			panic("ctx in null")
		}
		control.ctx.Shutdown()
	}()

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
	BotMode         bool
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
	Buffer                int
	shutdownChannel       chan struct{}
	receiveCh             chan string
	shellCtx              *ishell.Context
	multiChoiceCh         chan Message
	outputCh              chan Message
	workflow              string
	listWorflows          ListWorflows
	waitingResponse       bool
	disableInitialization bool
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
				bot.shellCtx.Println("Press Ctrl+C or type 'exit' to quit")
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
			if !message.BotMode {
				bot.waitingResponse = false
			}
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

func (bot *CliBot) SendMessage(message Message) error {

	if !bot.IsRunning() {
		return errors.New("Cli bot no running error")
	}

	if message.IsMultiChoice() {
		bot.multiChoiceCh <- message
		return nil
	}

	bot.outputCh <- message
	return nil
}
