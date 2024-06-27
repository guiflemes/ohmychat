package storage

import (
	"sync"

	"oh-my-chat/src/models"
)

type MemoryChatbotRepo struct {
	bots map[string]*models.ChatBot
	lock *sync.Mutex
}

func (m *MemoryChatbotRepo) GetChatBot(botName string) *models.ChatBot {
	m.lock.Lock()
	defer m.lock.Unlock()
	chatbot, ok := m.bots[botName]

	if !ok {
		return nil
	}
	return chatbot
}

func NewMemoryChatbotRepo() *MemoryChatbotRepo {
	return &MemoryChatbotRepo{
		bots: map[string]*models.ChatBot{
			"notion_notifierX_bot": {
				BotName:    "notion_notifierX_bot",
				Engine:     "guided",
				WorkflowID: "pokemon",
			},
			"cli_test": {
				BotName:    "cli_test",
				Engine:     "guided",
				WorkflowID: "pokemon",
			},
		},
		lock: &sync.Mutex{},
	}
}
