package core

import (
	"go.uber.org/zap"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/message"
	"oh-my-chat/src/models"
)

type Workflow interface {
	Engine() string
}

type WorkflowGetter interface {
	GetFlow(channelName string) Workflow
}

type Engine interface {
	Name() string
	Match(models.Message, <-chan models.Message)
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

type Matcher struct {
	workflowGetter WorkflowGetter
	workflow       Workflow
	engines        Engines
}

func (m *Matcher) Match(inputMsg <-chan models.Message, outputMsg chan<- models.Message) {
	for {
		message := <-inputMsg

		if m.workflow == nil {
			m.workflow = m.workflowGetter.GetFlow(message.ChannelName)
		}
	}
}

func (m *Matcher) handleWorkflow(msg models.Message, output chan<- models.Message) {
	if m.workflow == nil {
		logger.Logger.Error(
			"Work flow not found",
			zap.String("platfotm", string(msg.Remote)),
		)

		msg.Error = "some error has ocurred"
		output <- msg
		return
	}

	engine := m.engines.GetTarget(m.workflow.Engine())
	if engine == nil {
		logger.Logger.Error("Engine not found", zap.String("engine", m.workflow.Engine()))
		msg.Error = "some error has ocurred"
		output <- msg
		return
	}
}
