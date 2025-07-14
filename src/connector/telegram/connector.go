package telegram

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"oh-my-chat/src/core"
	"oh-my-chat/src/message"
	"oh-my-chat/src/utils"
)

type telegram struct {
	client *tgbotapi.BotAPI
}

func NewTelegramConnector(client *tgbotapi.BotAPI) core.Connector {
	return &telegram{client: client}
}

func (t *telegram) Acquire(ctx *core.ChatContext, input chan<- message.Message) error {

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

			msg := message.NewMessage()
			msg.Type = message.MsgTypeUnknown
			msg.Connector = message.Telegram
			msg.ConnectorID = strconv.Itoa(m.MessageID)
			msg.Input = m.Text
			msg.Service = message.MsgServiceChat
			msg.ChannelID = strconv.FormatInt(m.Chat.ID, 10)
			msg.BotID = strconv.FormatInt(user.ID, 10)
			msg.BotName = user.UserName

			input <- msg

		case <-ctx.Done():
			return nil
		}

	}
}
func (t *telegram) Dispatch(message message.Message) error {
	chatID, err := strconv.ParseInt(message.ChannelID, 10, 64)
	if err != nil {
		log.Printf("telegram: error parsing chat_id | %s", err)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, message.Output)
	t.formatResponse(&msg, message)

	_, err = t.client.Send(msg)
	if err != nil {
		log.Printf("telegram: error sending message '%s' | %s", message.ID, err)
		return err
	}
	return nil
}

func (t *telegram) formatResponse(responseMsg *tgbotapi.MessageConfig, msg message.Message) {
	switch msg.ResponseType {
	case message.OptionResponse:
		buttons := utils.Map(msg.Options, func(o message.Option) tgbotapi.InlineKeyboardButton {
			return tgbotapi.NewInlineKeyboardButtonData(o.Name, o.ID)
		})
		keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
		responseMsg.ReplyMarkup = keyboard
	default:
	}
}
