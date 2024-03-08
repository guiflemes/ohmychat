package telegram

import (
	"log"
	"notion-agenda/settings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// 6870062760
type WorkFlowGetter interface {
	Get(workflowId string) (WorkFlow, error)
}

type WorkFlow interface {
	Reply(string) (string, error)
}

type TelegramGate struct {
	client       *tgbotapi.BotAPI
	workFlowRepo WorkFlowGetter
}

func NewTelegramGate() *TelegramGate {

	client, err := tgbotapi.NewBotAPI(settings.GETENV("TELEGRAM_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	return &TelegramGate{client: client}
}

func (t *TelegramGate) SendDirectMessage(userID int64, message string) error {
	msg := tgbotapi.NewMessage(userID, message)
	_, err := t.client.Send(msg)

	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("message has been sent successfully")
	return nil
}

func (t *TelegramGate) Chat(userID int64, message string, timeout int, workflowId string) {
	err := t.SendDirectMessage(userID, message)
	if err != nil {
		log.Println(err)
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeout

	workflow, err := t.workFlowRepo.Get(workflowId)
	if err != nil {
		log.Println(err)
		return
	}

	updates := t.client.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		reply, err := workflow.Reply(update.Message.Text)

		if err != nil {
			reply = "I did not understand, can you repeat?"

		}

		replyMsg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
		replyMsg.ReplyToMessageID = update.Message.MessageID

		_, err = t.client.Send(replyMsg)
		if err != nil {
			log.Println(err)
		}

	}

}
