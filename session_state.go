package ohmychat

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type StateID string

type StateRegister struct {
	states map[StateID]SessionState
	mu     sync.Mutex
}

func NewStateRegister() *StateRegister {
	return &StateRegister{states: make(map[StateID]SessionState)}
}

func (s *StateRegister) GetState(stateID StateID) SessionState {
	s.mu.Lock()
	defer s.mu.Unlock()
	if state, ok := s.states[stateID]; ok {
		return state
	}
	return IdleState{}
}

func (s *StateRegister) Register(state SessionState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[state.Hash()] = state
}

type SessionState interface {
	IsState()
	Hash() StateID
}

type IdleState struct{}

func (IdleState) IsState() {}

func (IdleState) Hash() StateID {
	sum := sha256.Sum256([]byte("idle_state"))
	return StateID(sum[0:8])
}

type WaitingInputState struct {
	PromptEmptyMessage string
	PromptExit         string
	ExitInput          string // do not use exit as input for cli connector is a reserved keyword for it
	Action             ActionFunc
}

func (WaitingInputState) IsState() {}

func (w WaitingInputState) Hash() StateID {
	handlerPtr := fmt.Sprintf("%p", w.Action)
	data := fmt.Sprintf("%s%s%s%s", w.PromptEmptyMessage, w.PromptExit, w.ExitInput, handlerPtr)
	sum := sha256.Sum256([]byte(data))
	return StateID(sum[0:8])
}

type WaitingChoiceState struct {
	Prompt              string
	PromptInvalidOption string
	Choices             Choices
}

func (WaitingChoiceState) IsState() {}

func (w WaitingChoiceState) Hash() StateID {
	data := fmt.Sprintf("%s%s", w.Prompt, w.PromptInvalidOption)
	for _, ch := range w.Choices {
		handlerPtr := fmt.Sprintf("%p", ch.Handler)
		data += fmt.Sprintf("%s%s", ch.Text, handlerPtr)
	}
	sum := sha256.Sum256([]byte(data))
	return StateID(fmt.Sprintf("%x", sum[:8]))
}

type Choice struct {
	Text    string
	Handler ActionFunc
}

type Choices []Choice

func (c Choices) GetHandler(option string) (ActionFunc, bool) {
	for _, opt := range c {
		if opt.Text == option {
			return opt.Handler, true
		}
	}
	return nil, false
}

func (c Choices) BindMany(handler ActionFunc, options ...string) Choices {
	for _, opt := range options {
		c = append(c, Choice{Text: opt, Handler: handler})
	}
	return c
}

type WaitingBotResponseState struct {
	OnDone ActionFunc
}

func (WaitingBotResponseState) IsState() {}

func (s WaitingBotResponseState) Hash() StateID {
	handlerPtr := fmt.Sprintf("%p", s.OnDone)
	data := fmt.Sprintf("waiting_bot_response%s", handlerPtr)
	sum := sha256.Sum256([]byte(data))
	return StateID(fmt.Sprintf("%x", sum[:8]))
}
