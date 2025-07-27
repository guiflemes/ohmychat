package connector

import (
	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/connector/cli"
	"github.com/guiflemes/ohmychat/connector/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CliOption = cli.CliOption

func Cli(options ...CliOption) ohmychat.Connector {
	return cli.NewCliConnector(options...)
}

func Telegram(tgbot *tgbotapi.BotAPI) ohmychat.Connector {
	return telegram.NewTelegramConnector(tgbot)
}
