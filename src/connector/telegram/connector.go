package telegram

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"oh-my-chat/src/connector"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
	"oh-my-chat/src/utils"
)

type telegram struct {
	client *tgbotapi.BotAPI
}

func NewTelegramConnector(bot *models.Bot) (connector.Connector, error) {
	client, err := tgbotapi.NewBotAPI(bot.TelegramConfig.Token)

	if err != nil {
		return nil, err
	}

	return &telegram{client: client}, nil

}

func (t *telegram) Acquire(ctx context.Context, input chan<- models.Message) {

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

	fmt.Println("botuser", user) // TODO -> put on metadata

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

			message := models.NewMessage()
			message.Type = models.MsgTypeUnknown
			message.Connector = models.Telegram
			message.ConnectorID = strconv.Itoa(m.MessageID)
			message.Input = m.Text
			message.Service = models.MsgServiceChat
			message.ChannelID = strconv.FormatInt(m.Chat.ID, 10)
			message.BotID = strconv.FormatInt(user.ID, 10)
			message.BotName = user.UserName

			input <- message

		case <-ctx.Done():
			logger.Logger.Info(
				"context cancelled, stopping Acquire",
				zap.String("context", "telegram_client"),
			)
			return
		}

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
	t.formatResponse(&msg, message)

	_, err := t.client.Send(msg)
	if err != nil {
		logger.Logger.Error("unable to send message", zap.Error(err),
			zap.String("context", "telegram_client"))
	}
}

func (t *telegram) formatResponse(responseMsg *tgbotapi.MessageConfig, message models.Message) {
	switch message.ResponseType {
	case models.OptionResponse:
		buttons := utils.Map(message.Options, func(o models.Option) tgbotapi.InlineKeyboardButton {
			return tgbotapi.NewInlineKeyboardButtonData(o.Name, o.ID)
		})
		keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
		responseMsg.ReplyMarkup = keyboard
	default:
	}
}
