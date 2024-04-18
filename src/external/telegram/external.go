package telegram

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

type telegram struct {
	client *tgbotapi.BotAPI
}

func (t *telegram) Acquire(input chan<- models.Message) {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	user, err := t.client.GetMe()
	if err != nil {
		logger.Logger.Error(
			"failed to initialize telegram client",
			zap.Error(err),
			zap.String("context", "telegram_client"),
		)
		return
	}

	fmt.Println("botuser", user)

	updates := t.client.GetUpdatesChan(u)

	for update := range updates {

		var m *tgbotapi.Message

		if update.Message != nil {
			m = update.Message
		}

		if update.ChannelPost != nil {
			m = update.ChannelPost
		}

		if m == nil {
			continue
		}

		if m.From != nil && m.From.ID == user.ID {
			continue
		}

		message := models.NewMessage()
		message.Type = models.MsgTypeUnknown
		message.Remote = models.Telegram
		message.RemoteID = strconv.Itoa(m.MessageID)
		message.Input = m.Text
		message.Service = models.MsgServiceChat
		message.ChannelID = strconv.FormatInt(m.Chat.ID, 10)

		input <- message
	}
}
func (t *telegram) Dispatch(message models.Message) {
	chatID, error := strconv.ParseInt(message.ChannelID, 10, 64)
	if error != nil {
		logger.Logger.Error(
			"unable to retrieve chat",
			zap.Error(error),
			zap.Int64("chat_id", chatID),
			zap.String("context", "telegram_client"),
		)
		return
	}

	msg := tgbotapi.NewMessage(chatID, message.Output)

	_, err := t.client.Send(msg)
	if err != nil {
		logger.Logger.Error("unable to send message", zap.Error(err),
			zap.String("context", "telegram_client"))
	}
}
