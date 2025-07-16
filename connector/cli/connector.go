package cli

import (
	"fmt"

	"github.com/abiosoft/ishell"
	"github.com/guiflemes/ohmychat/utils"

	"github.com/guiflemes/ohmychat"
)

type BotCli interface {
	GetUpdateChanels() UpdateChannel
	SendMessage(message Message) error
}

type ChatControl struct {
	ctx *ohmychat.ChatContext
}

type cliConnector struct {
	bot     BotCli
	control *ChatControl
}

func NewCliConnector(options ...CliOption) ohmychat.Connector {
	control := &ChatControl{}
	cliBot := NewCliBot(ishell.New(), control, options...)
	conn := &cliConnector{bot: cliBot, control: control}
	return conn
}

func (cli *cliConnector) Acquire(ctx *ohmychat.ChatContext, input chan<- ohmychat.Message) error {
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
			msg := ohmychat.NewMessage()
			msg.Type = ohmychat.MsgTypeUnknown
			msg.Connector = ohmychat.Cli
			msg.ConnectorID = "CLI"
			msg.Input = update.Message.Text
			msg.Service = ohmychat.MsgServiceChat
			msg.ChannelID = "CLI"
			msg.BotID = "CLI"
			msg.BotName = update.Message.BotName
			msg.User.ID = "cli_id"

			input <- msg

		default:
		}
	}

}

func (cli *cliConnector) Dispatch(msg ohmychat.Message) error {
	resposeMsg := NewMessage(msg.Output)
	resposeMsg.UnBlockByAction = msg.ActionDone

	switch msg.ResponseType {
	case ohmychat.OptionResponse:
		options := utils.Map(msg.Options, func(o ohmychat.Option) string {
			return o.ID
		})
		resposeMsg.MultiChoice = options
	default:
	}

	return cli.bot.SendMessage(resposeMsg)
}
