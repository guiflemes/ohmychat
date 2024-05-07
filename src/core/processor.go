package core

import (
	"context"

	"go.uber.org/zap"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

type Workflow interface {
	Engine() string
}

type ChatBotGetter interface {
	GetChatBot(botName string) *models.ChatBot
}

type ActionQueue interface {
	Put(actionPair ActionReplyPair)
}

type Engine interface {
	Name() string
	HandleMessage(models.Message, chan<- models.Message)
	GetActionQueue() ActionQueue
	Config(channelName string) error
	IsReady() bool
}

type Engines []Engine

func (e Engines) GetTarget(target string) Engine {
	for _, eng := range e {
		if eng.Name() == target {
			return eng
		}
	}
	return nil
}

type processor struct {
	chatBotGetter ChatBotGetter
	chatBot       *models.ChatBot
	engines       Engines
}

func NewProcessor(chatBotGetter ChatBotGetter, engines Engines) *processor {
	return &processor{
		chatBotGetter: chatBotGetter,
		engines:       engines,
	}
}

func (m *processor) Process(
	ctx context.Context,
	inputMsg <-chan models.Message,
	outputMsg chan<- models.Message,
) {
	for {
		select {
		case message, ok := <-inputMsg:
			if !ok {
				return
			}
			if m.chatBot == nil {
				m.chatBot = m.chatBotGetter.GetChatBot(message.BotName)
			}

			m.handleWorkflow(message, outputMsg)

		case <-ctx.Done():
			return
		}

	}
}

func (m *processor) handleWorkflow(msg models.Message, output chan<- models.Message) {
	if m.chatBot == nil {
		logger.Logger.Error(
			"Chatbot not found",
			zap.String(
				"platfotm",
				string(msg.Connector),
			),
			zap.String("context", "processor"),
			zap.String("chatbot", msg.BotName),
		)
		response := &msg
		response.Output = "some error has ocurred"
		response.Error = "chat not found"
		output <- *response
		return
	}

	engine := m.engines.GetTarget(m.chatBot.Engine)

	if engine == nil {
		logger.Logger.Error(
			"Engine not found",
			zap.String("engine", m.chatBot.Engine), zap.String("context", "processor"),
		)
		response := &msg
		response.Output = "some error has ocurred"
		response.Error = "engine not found"
		output <- *response
		return
	}

	if !engine.IsReady() {
		logger.Logger.Info(
			"The engine is not ready. Starting to prepare it",
			zap.String("engine", engine.Name()),
		)
		engine.Config(m.chatBot.WorkflowID)
	}

	logger.Logger.Info("handling message", zap.String("message_id", msg.BotID))

	engine.HandleMessage(msg, output)
}
