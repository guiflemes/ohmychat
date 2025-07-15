package cli

import (
	"fmt"

	"github.com/abiosoft/ishell"

	"github.com/guiflemes/ohmychat/src/message"
	"github.com/guiflemes/ohmychat/src/utils"

	"github.com/guiflemes/ohmychat/src/core"
)

type BotCli interface {
	GetUpdateChanels() UpdateChannel
	SendMessage(message Message) error
}

type ChatControl struct {
	ctx *core.ChatContext
}

type cliConnector struct {
	bot     BotCli
	control *ChatControl
}

func NewCliConnector(options ...CliOption) core.Connector {
	control := &ChatControl{}
	cliBot := NewCliBot(ishell.New(), control, options...)
	conn := &cliConnector{bot: cliBot, control: control}
	return conn
}

func (cli *cliConnector) Acquire(ctx *core.ChatContext, input chan<- message.Message) error {
	cli.control.ctx = ctx

	updates := cli.bot.GetUpdateChanels()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("sutdown cli connector")
			return nil
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

func (cli *cliConnector) Dispatch(msg message.Message) error {
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

	return cli.bot.SendMessage(resposeMsg)
}
