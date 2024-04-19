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

type WorkflowGetter interface {
	GetFlow(channelName string) Workflow
}

type ActionQueue interface {
	Consume(context.Context)
	Put(actionPair ActionReplyPair)
}

type Engine interface {
	Name() string
	HandleMessage(models.Message, chan<- models.Message)
	GetActionQueue() ActionQueue
	Config(workflow Workflow)
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
	workflowGetter WorkflowGetter
	workflow       Workflow
	engines        Engines
}

func NewProcessor(workflowGetter WorkflowGetter, engines Engines) *processor {
	return &processor{
		workflowGetter: workflowGetter,
		engines:        engines,
	}
}

func (m *processor) Process(inputMsg <-chan models.Message, outputMsg chan<- models.Message) {
	for {
		message := <-inputMsg

		if m.workflow == nil {
			m.workflow = m.workflowGetter.GetFlow(message.ChannelName)
		}

		m.handleWorkflow(message, outputMsg)
	}
}

func (m *processor) handleWorkflow(msg models.Message, output chan<- models.Message) {
	if m.workflow == nil {
		logger.Logger.Error(
			"Work flow not found",
			zap.String("platfotm", string(msg.Connector)), zap.String("context", "processor"),
		)

		msg.Error = "some error has ocurred"
		output <- msg
		return
	}

	engine := m.engines.GetTarget(m.workflow.Engine())

	if engine == nil {
		logger.Logger.Error(
			"Engine not found",
			zap.String("engine", m.workflow.Engine()), zap.String("context", "processor"),
		)
		msg.Error = "some error has ocurred"
		output <- msg
		return
	}

	if !engine.IsReady() {
		logger.Logger.Info(
			"The engine is not ready. Starting to prepare it",
			zap.String("engine", engine.Name()),
		)
		engine.Config(m.workflow)
	}

	engine.HandleMessage(msg, output)
}
