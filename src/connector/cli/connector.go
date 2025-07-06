package cli

import (
	"fmt"

	"github.com/abiosoft/ishell"

	"oh-my-chat/src/connector"
	"oh-my-chat/src/message"
	"oh-my-chat/src/utils"

	"oh-my-chat/src/context"
)

type BotCli interface {
	GetUpdateChanels() UpdateChannel
	SendMessage(message Message)
}

type ChatControl struct {
	ctx *context.ChatContext
}

type cliConnector struct {
	bot     BotCli
	control *ChatControl
}

func NewCliConnector(options ...CliOption) connector.Connector {
	control := &ChatControl{}
	cliBot := NewCliBot(ishell.New(), control, options...)
	conn := &cliConnector{bot: cliBot, control: control}
	return conn
}

func (cli *cliConnector) Acquire(ctx *context.ChatContext, input chan<- message.Message) {
	cli.control.ctx = ctx

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
			msg := message.NewMessage()
			msg.Type = message.MsgTypeUnknown
			msg.Connector = message.Cli
			msg.ConnectorID = "CLI"
			msg.Input = update.Message.Text
			msg.Service = message.MsgServiceChat
			msg.ChannelID = "CLI"
			msg.BotID = "CLI"
			msg.BotName = update.Message.BotName
			msg.User.ID = "cli_id"

			input <- msg

		default:
		}
	}

}

func (cli *cliConnector) Dispatch(msg message.Message) {
	resposeMsg := NewMessage(msg.Output)
	resposeMsg.UnBlockByAction = msg.ActionDone

	switch msg.ResponseType {
	case message.OptionResponse:
		options := utils.Map(msg.Options, func(o message.Option) string {
			return o.ID
		})
		resposeMsg.MultiChoice = options
	default:
	}

	cli.bot.SendMessage(resposeMsg)

}
