package telegram

import (
	"log"
	"notion-agenda/src/notion"
	"notion-agenda/src/service"
)

type PendencyFormatter interface {
	Format(pendency []notion.StudyStep) (string, error)
}
type TelegramMsgSender interface {
	SendDirectMessage(userID int64, message string) error
}

type telegramPendencyHandler struct {
	formatter PendencyFormatter
	sender    TelegramMsgSender
}

func NewTelegramPendencyHandler() *telegramPendencyHandler {
	return &telegramPendencyHandler{
		formatter: &pendencyFormatter{},
		sender:    NewTelegramGate(),
	}
}

func (h *telegramPendencyHandler) Handle(message service.Message) error {
	pendencies, ok := message.(*notion.PendencyEvent)

	if !ok {
		log.Printf("Unexpected type in Function: %T", message)
		panic("Critical error: Unexpected type")
	}

	msg, err := h.formatter.Format(pendencies.Pendency)

	if err != nil {
		return err
	}

	return h.sender.SendDirectMessage(6870062760, msg)
}
