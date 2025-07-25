//go:generate mockgen -source engine.go -destination ./mocks/engine.go -package mocks
package rule_engine

import (
	"strings"
	"time"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/utils"
)

const MatchAll = "__all__"

type Rule struct {
	Prompts   []string
	Action    ohmychat.ActionFunc
	NextState ohmychat.SessionState
}

type MatcherFunc func(rules []Rule, input string) (Rule, bool)
type RuleEngineOption func(engine *RuleEngine)

func WithMatcher(m MatcherFunc) RuleEngineOption {
	return func(engine *RuleEngine) {
		engine.matcher = m
	}
}

func WithSessionExpiresAt(s time.Duration) RuleEngineOption {
	return func(engine *RuleEngine) {
		engine.sessionExpiresAt = utils.PtrOf(s)
	}
}

type RuleEngine struct {
	matcher          MatcherFunc
	rules            []Rule
	sessionExpiresAt *time.Duration
}

func NewRuleEngine(opts ...RuleEngineOption) *RuleEngine {
	engine := &RuleEngine{}

	for _, opt := range opts {
		opt(engine)
	}

	if engine.matcher == nil {
		engine.matcher = DefaultMatcher
	}

	if engine.sessionExpiresAt == nil {
		engine.sessionExpiresAt = utils.PtrOf(ohmychat.SessionExpiresAt)
	}

	return engine
}

func (e *RuleEngine) RegisterRule(rule ...Rule) {
	e.rules = append(e.rules, rule...)
}

func (e *RuleEngine) HandleMessage(ctx *ohmychat.Context, msg *ohmychat.Message) {
	sess := ctx.Session()

	if sess.IsExpired(*e.sessionExpiresAt) {
		sess.State = ohmychat.IdleState{}
	}

	switch state := sess.State.(type) {
	case ohmychat.IdleState:
		e.handleIdleState(ctx, msg)
	case ohmychat.WaitingInputState:
		e.handleWaitingInputState(ctx, msg, state)
	case ohmychat.WaitingChoiceState:
		e.handleWaitingChoiceState(ctx, msg, state)
	default:
		e.handleUnknownState(ctx, msg)
	}
}

func (e *RuleEngine) handleIdleState(ctx *ohmychat.Context, msg *ohmychat.Message) {
	rule, ok := e.matcher(e.rules, msg.Input)
	if !ok {
		msg.Output = "desculpe não entendi"
		ctx.SendOutput(msg)
		return
	}

	ctx.SetSessionState(rule.NextState)
	rule.Action(ctx, msg)

}

func (e *RuleEngine) handleWaitingInputState(ctx *ohmychat.Context, msg *ohmychat.Message, state ohmychat.WaitingInputState) {
	if strings.TrimSpace(msg.Input) == "" {
		msg.Output = state.PromptEmptyMessage
		ctx.SendOutput(msg)
		return
	}
	if msg.Input == state.ExitInput {
		ctx.SetSessionState(ohmychat.IdleState{})
		msg.Output = state.PromptExit
		ctx.SendOutput(msg)
		return
	}
	state.Action(ctx, msg)
}

func (e *RuleEngine) handleWaitingChoiceState(ctx *ohmychat.Context, msg *ohmychat.Message, state ohmychat.WaitingChoiceState) {
	handler, ok := state.Choices[msg.Input]
	if !ok {
		msg.Output = state.PromptInvalidOption
		ctx.SendOutput(msg)
		return
	}

	ctx.SetSessionState(ohmychat.IdleState{})
	handler(ctx, msg)

}

func (e *RuleEngine) handleUnknownState(ctx *ohmychat.Context, msg *ohmychat.Message) {
	msg.Output = "Erro interno: estado desconhecido."
	ctx.SendOutput(msg)
}

func matchInsensitiveContains(input, pattern string) bool {
	return strings.Contains(strings.ToLower(input), strings.ToLower(pattern))
}

func DefaultMatcher(rules []Rule, input string) (Rule, bool) {
	for _, rule := range rules {
		for _, pattern := range rule.Prompts {
			if pattern == MatchAll {
				return rule, true
			}
			if matchInsensitiveContains(input, pattern) {
				return rule, true
			}
		}
	}
	return Rule{}, false
}
