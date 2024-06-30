package chat

import (
	"errors"

	"oh-my-chat/src/config"
	"oh-my-chat/src/models"
)

type ChatDatabase interface {
	GetChatBot(botName string) *models.ChatBot
	ListChatBots() []*models.ChatBot
}

func NewConnection(config config.ChatDatabase) (ChatDatabase, error) {
	switch config.Kind {
	case "memory":
		return NewMemoryChatbotRepo(), nil
	}

	return nil, errors.New("unable to connect database")
}
