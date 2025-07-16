package telegram

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/utils"
)

type telegram struct {
	client *tgbotapi.BotAPI
}

func NewTelegramConnector(client *tgbotapi.BotAPI) ohmychat.Connector {
	return &telegram{client: client}
}

func (t *telegram) Acquire(ctx *ohmychat.ChatContext, input chan<- ohmychat.Message) error {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	user, err := t.client.GetMe()
	if err != nil {
		return err
	}

	updates := t.client.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:

			var m *tgbotapi.Message

			if update.Message != nil {
				m = update.Message
			}

			if update.ChannelPost != nil {
				m = update.ChannelPost
			}

			if update.CallbackQuery != nil {
				m = update.CallbackQuery.Message
				m.Text = update.CallbackData()
			}

			if m == nil {
				continue
			}

			msg := ohmychat.NewMessage()
			msg.Type = ohmychat.MsgTypeUnknown
			msg.Connector = ohmychat.Telegram
			msg.ConnectorID = strconv.Itoa(m.MessageID)
			msg.Input = m.Text
			msg.Service = ohmychat.MsgServiceChat
			msg.ChannelID = strconv.FormatInt(m.Chat.ID, 10)
			msg.BotID = strconv.FormatInt(user.ID, 10)
			msg.BotName = user.UserName

			input <- msg

		case <-ctx.Done():
			return nil
		}

	}
}
func (t *telegram) Dispatch(ohmychat ohmychat.Message) error {
	chatID, err := strconv.ParseInt(ohmychat.ChannelID, 10, 64)
	if err != nil {
		log.Printf("telegram: error parsing chat_id | %s", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, ohmychat.Output)
	t.formatResponse(&msg, ohmychat)

	_, err = t.client.Send(msg)
	if err != nil {
		log.Printf("telegram: error sending ohmychat '%s' | %s", ohmychat.ID, err)
		return err
	}
	return nil
}

func (t *telegram) formatResponse(responseMsg *tgbotapi.MessageConfig, msg ohmychat.Message) {
	switch msg.ResponseType {
	case ohmychat.OptionResponse:
		buttons := utils.Map(msg.Options, func(o ohmychat.Option) tgbotapi.InlineKeyboardButton {
			return tgbotapi.NewInlineKeyboardButtonData(o.Name, o.ID)
		})
		keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
		responseMsg.ReplyMarkup = keyboard
	default:
	}
}
