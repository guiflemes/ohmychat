package rule_engine

type SessionState interface {
	IsState()
}

type IdleState struct{}

func (IdleState) IsState() {}

type WaitingInputState struct {
	Prompt string
	Action ActionFunc
}

func (WaitingInputState) IsState() {}

type WaitingChoiceState struct {
	Prompt  string
	Choices map[string]ActionFunc
}

func (WaitingChoiceState) IsState() {}
