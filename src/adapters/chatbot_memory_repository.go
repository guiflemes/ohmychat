package adapters

import (
	"sync"

	"oh-my-chat/src/models"
)

type MemoryChatbotRepo struct {
	bots map[string]*models.ChatBot
	lock *sync.Mutex
}

func (m *MemoryChatbotRepo) GetChatBot(channelName string) *models.ChatBot {
	m.lock.Lock()
	defer m.lock.Unlock()
	chatbot, ok := m.bots["my_first_bot"]

	if !ok {
		return nil
	}
	return chatbot
}

func NewMemoryChatbotRepo() *MemoryChatbotRepo {
	return &MemoryChatbotRepo{
		bots: map[string]*models.ChatBot{
			"my_first_bot": {
				BotName:    "my_first_bot",
				Engine:     "guided",
				WorkflowID: "pokemon",
			},
		},
		lock: &sync.Mutex{},
	}
}
