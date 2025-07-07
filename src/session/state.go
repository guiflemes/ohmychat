package session

import (
	"context"
	"oh-my-chat/src/message"
)

type ActionInput struct {
	Session *Session
	Message *message.Message
	Output  chan<- message.Message
}

type ActionFunc func(ctx context.Context, input ActionInput)

type SessionState interface {
	IsState()
}

type IdleState struct{}

func (IdleState) IsState() {}

type WaitingInputState struct {
	PromptEmptyMessage string
	Action             ActionFunc
}

func (WaitingInputState) IsState() {}

type WaitingChoiceState struct {
	Prompt  string
	Choices Choices
}

func (WaitingChoiceState) IsState() {}

type Choices map[string]ActionFunc
