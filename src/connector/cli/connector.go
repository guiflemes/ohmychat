package cli

import (
	"context"
	"fmt"

	"github.com/abiosoft/ishell"

	"oh-my-chat/src/connector"
	"oh-my-chat/src/models"
	"oh-my-chat/src/utils"
)

type BotCli interface {
	GetUpdateChanels() UpdateChannel
	SendMessage(message Message)
}

type cliConnector struct {
	bot BotCli
}

func NewCliConnector(bot *models.Bot) (connector.Connector, error) {
	cliBot := NewCliBot(bot, ishell.New())
	conn := &cliConnector{bot: cliBot}
	return conn, nil
}

func (cli *cliConnector) Acquire(ctx context.Context, input chan<- models.Message) {

	updates := cli.bot.GetUpdateChanels()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("sutdown cli connector")
			return
		case update, ok := <-updates:
			if !ok {
				continue
			}
			message := models.NewMessage()
			message.Type = models.MsgTypeUnknown
			message.Connector = models.Cli
			message.ConnectorID = "CLI"
			message.Input = update.Message.Text
			message.Service = models.MsgServiceChat
			message.ChannelID = "CLI"
			message.BotID = "CLI"
			message.BotName = update.Message.BotName
			message.User.ID = "cli_id"

			input <- message

		default:
		}
	}

}

func (cli *cliConnector) Dispatch(message models.Message) {
	resposeMsg := NewMessage(message.Output)
	resposeMsg.UnBlockByAction = message.ActionDone

	switch message.ResponseType {
	case models.OptionResponse:
		options := utils.Map(message.Options, func(o models.Option) string {
			return o.ID
		})
		resposeMsg.MultiChoice = options
	default:
	}

	cli.bot.SendMessage(resposeMsg)

}
