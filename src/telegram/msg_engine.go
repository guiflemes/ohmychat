package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"

	"oh-my-chat/settings"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/message"
	"oh-my-chat/src/utils"
)

type commandEngine struct {
	tree         *message.MessageTree
	node         *message.MessageNode
	dialogLaunch bool
	actionQueue  ActionQueue
}

func (e *commandEngine) IsInitialized() bool {
	return e.tree != nil || e.node != nil
}

func (e *commandEngine) resolveMessageNode(messageID string) {

	if e.dialogLaunch {
		node := e.node.SearchOneLevel(messageID)

		if node == nil {
			e.node = e.tree.Root()
			e.dialogLaunch = false
			return
		}
		e.node = node
		return
	}

	e.dialogLaunch = true

}

func (e *commandEngine) Reply(
	ctx context.Context,
	chatID int64,
	messageID string,
) tgbotapi.MessageConfig {
	e.resolveMessageNode(messageID)

	if e.node.Message().HasAction() {
		go func() {
			action := e.node.Message().Action
			content := action.Execute(ctx, messageID)
			e.actionQueue <- tgbotapi.NewMessage(chatID, content)
		}()
	}

	buttons := make([]tgbotapi.KeyboardButton, 0)

	e.node.TransverseInChildren(func(child *message.MessageNode) {
		buttons = append(buttons, tgbotapi.NewKeyboardButton(child.Message().ID()))
	})

	keyboard := tgbotapi.NewReplyKeyboard(buttons)

	msg := tgbotapi.NewMessage(chatID, e.node.Message().Content)
	msg.ReplyMarkup = keyboard
	return msg

}

type ActionQueue chan tgbotapi.MessageConfig

type WorkFlowEngine struct {
	client           *tgbotapi.BotAPI
	notRecognizedMsg string
	dialogLaunch     bool
	commandEngine    *commandEngine
	actionQueue      ActionQueue
}

func NewEngine() *WorkFlowEngine {
	client, err := tgbotapi.NewBotAPI(settings.GETENV("TELEGRAM_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	commandTree := message.PokemonFlow()
	actionQueue := make(chan tgbotapi.MessageConfig)

	return &WorkFlowEngine{
		client:           client,
		notRecognizedMsg: "I did not understand, can you repeat?",
		dialogLaunch:     false,
		commandEngine: &commandEngine{
			tree:        commandTree,
			node:        commandTree.Root(),
			actionQueue: actionQueue,
		},
		actionQueue: actionQueue,
	}
}

func (e *WorkFlowEngine) HasPostback() bool {
	return e.commandEngine.IsInitialized()
}

func (e *WorkFlowEngine) Chating(timeout int) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	updates := e.client.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			e.replyMessage(update)

		case action := <-e.actionQueue:
			_, err := e.client.Send(action)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

func (e *WorkFlowEngine) replyMessage(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	var user *utils.User
	sentFrom := update.SentFrom()
	if sentFrom != nil {
		user = &utils.User{
			ID:        strconv.FormatInt(sentFrom.ID, 10),
			Provider:  "telegram",
			FirstName: sentFrom.FirstName,
			LastName:  sentFrom.LastName,
			Username:  sentFrom.UserName,
		}
	}

	ctx := context.WithValue(context.Background(), utils.UserKey, user)

	logger.Logger.Info(
		"Message received",
		zap.String("platfotm", "telegram"),
		zap.String("user_id", user.ID),
		zap.String("user_name", fmt.Sprintf("%s %s", user.FirstName, user.LastName)),
	)

	replyMsg := e.commandEngine.Reply(ctx, update.Message.Chat.ID, update.Message.Text)
	replyMsg.ReplyToMessageID = update.Message.MessageID

	_, err := e.client.Send(replyMsg)
	if err != nil {
		log.Println(err)
	}

}
