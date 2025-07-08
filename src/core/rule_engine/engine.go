package rule_engine

import (
	"oh-my-chat/src/core"
	"oh-my-chat/src/message"
	"strings"
)

type Rule struct {
	Prompts   []string
	Action    core.ActionFunc
	NextState core.SessionState
}

type MatcherFunc func(rules []Rule, input string) (Rule, bool)
type RuleEngineOption func(engine *RuleEngine)

func WithMatcher(m MatcherFunc) RuleEngineOption {
	return func(engine *RuleEngine) {
		engine.matcher = m
	}
}

type RuleEngine struct {
	matcher MatcherFunc
	rules   []Rule
}

func NewRuleEngine(opts ...RuleEngineOption) *RuleEngine {
	engine := &RuleEngine{}

	for _, opt := range opts {
		opt(engine)
	}

	if engine.matcher == nil {
		engine.matcher = DefaultMatcher
	}

	return engine
}

func (e *RuleEngine) RegisterRule(rule ...Rule) {
	e.rules = append(e.rules, rule...)
}

func (e *RuleEngine) HandleMessage(ctx *core.Context, msg *message.Message) {
	sess := ctx.Session()

	switch state := sess.State.(type) {
	case core.IdleState:
		e.handleIdleState(ctx, msg)
	case core.WaitingInputState:
		e.handleWaitingInputState(ctx, msg, state)
	case core.WaitingChoiceState:
		e.handleWaitingChoiceState(ctx, msg, state)
	default:
		e.handleUnknownState(ctx, msg)
	}
}

func (e *RuleEngine) handleIdleState(ctx *core.Context, msg *message.Message) {
	rule, ok := e.matcher(e.rules, msg.Input)
	if !ok {
		msg.Output = "desculpe não entendi"
		ctx.SendOutput(msg)
		return
	}

	ctx.SetSessionState(rule.NextState)
	rule.Action(ctx, msg)

}

func (e *RuleEngine) handleWaitingInputState(ctx *core.Context, msg *message.Message, state core.WaitingInputState) {
	if strings.TrimSpace(msg.Input) == "" {
		msg.Output = state.PromptEmptyMessage
		ctx.SendOutput(msg)
		return
	}
	state.Action(ctx, msg)
}

func (e *RuleEngine) handleWaitingChoiceState(ctx *core.Context, msg *message.Message, state core.WaitingChoiceState) {
	handler, ok := state.Choices[msg.Input]
	if !ok {
		msg.Output = "Opção inválida. " + state.Prompt
		ctx.SendOutput(msg)
		return
	}

	ctx.SetSessionState(core.IdleState{})
	handler(ctx, msg)

}

func (e *RuleEngine) handleUnknownState(ctx *core.Context, msg *message.Message) {
	msg.Output = "Erro interno: estado desconhecido."
	ctx.SendOutput(msg)
}

func DefaultMatcher(rules []Rule, input string) (Rule, bool) {
	for _, rule := range rules {
		for _, pattern := range rule.Prompts {
			if strings.Contains(strings.ToLower(input), strings.ToLower(pattern)) {
				return rule, true
			}
		}
	}
	return Rule{}, false
}
