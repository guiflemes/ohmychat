package rule_engine

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
