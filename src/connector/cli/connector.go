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
}

type cliConnector struct {
	bot *CliBot
}

func NewCliConnector(bot *models.Bot) (connector.Connector, error) {
	shell := ishell.New()

	go func() { shell.Run() }()

	cliBot := NewCliBot()
	conn := &cliConnector{bot: cliBot}

	shell.AddCmd(&ishell.Cmd{
		Name: "chat",
		Help: "Marvin",
		Func: func(c *ishell.Context) {
			cliBot.StartChat(c)
		},
	})

	return conn, nil
}

func (cli *cliConnector) Acquire(ctx context.Context, input chan<- models.Message) {

	updates := cli.bot.GetUpdateChanels()

	for {
		select {
		case <-ctx.Done():
			cli.bot.StopReceivingUpdates()
			fmt.Println("sutdown shell")
			return
		case update := <-updates:
			message := models.NewMessage()
			message.Type = models.MsgTypeUnknown
			message.Connector = models.Cli
			message.ConnectorID = ""
			message.Input = update.Message.Text
			message.Service = models.MsgServiceChat
			message.ChannelID = ""
			message.BotID = ""
			message.BotName = "bot"

			//TODO used to text remove this block and Dispatch
			if update.Message.Text == "text" {
				message.ResponseType = models.TextResponse
			}
			cli.Dispatch(message)

		default:
		}
	}

}

func (cli *cliConnector) Dispatch(message models.Message) {
	resposeMsg := NewMessage(message.Input)

	switch message.ResponseType {
	case models.OptionResponse:
		options := utils.Map(message.Options, func(o models.Option) string {
			return o.Name
		})
		options = append(options, "text")
		options = append(options, "test2")
		options = append(options, "test3")
		resposeMsg.MultiChoice = options
	default:
	}

	cli.bot.SendMessage(resposeMsg)

}
