package core

type SessionState interface {
	IsState()
}

type IdleState struct{}

func (IdleState) IsState() {}

type WaitingInputState struct {
	PromptEmptyMessage string
	PromptExit         string
	ExitInput          string // do not use exit as input for cli connector is a reserved keyword for it
	Action             ActionFunc
}

func (WaitingInputState) IsState() {}

type WaitingChoiceState struct {
	Prompt              string
	PromptInvalidOption string
	Choices             Choices
}

func (WaitingChoiceState) IsState() {}

type Choices map[string]ActionFunc
