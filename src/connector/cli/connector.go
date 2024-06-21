package cli

import (
	"context"
	"fmt"

	"github.com/abiosoft/ishell"

	"oh-my-chat/src/connector"
	"oh-my-chat/src/models"
	"oh-my-chat/src/utils"
)

type CliMessage struct {
	Text string
	ctx  *ishell.Context
}

type cliConnector struct {
	shell *ishell.Shell
	msgCh chan CliMessage
}

func NewCliConnector(bot *models.Bot) (connector.Connector, error) {
	shell := ishell.New()
	msgCh := make(chan CliMessage, 1)

	go func() { shell.Run() }()

	shell.AddCmd(&ishell.Cmd{
		Name: "chat",
		Help: "Marvin",
		Func: func(c *ishell.Context) {
			startChat(c, msgCh)
		},
	})
	return &cliConnector{
		shell: shell,
		msgCh: msgCh,
	}, nil
}

func startChat(c *ishell.Context, msgCh chan<- CliMessage) {
	for {
		input := c.ReadLine()
		if input == "exit" {
			c.Println("Exiting chat mode...")
			break
		}

		msgCh <- CliMessage{Text: input, ctx: c}
	}
}

func (cli *cliConnector) Acquire(ctx context.Context, input chan<- models.Message) {

	for {
		select {
		case <-ctx.Done():
			cli.shell.Close()
			fmt.Println("sutdown shell")
			return
		case msg := <-cli.msgCh:
			message := models.NewMessage()
			message.Type = models.MsgTypeUnknown
			message.Connector = models.Cli
			message.ConnectorID = ""
			message.Input = msg.Text
			message.Service = models.MsgServiceChat
			message.ChannelID = ""
			message.BotID = ""
			message.BotName = "bot"

			cli.shell.Println("input ", msg.Text)
		default:
		}
	}

}

func (cli *cliConnector) Dispatch(message models.Message) {
	switch message.ResponseType {
	case models.OptionResponse:
		options := utils.Map(message.Options, func(o models.Option) string {
			return o.Name
		})
		cli.shell.Println(options)

	default:
	}

}
