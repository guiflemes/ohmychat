package rule_engine

import (
	"context"
	"oh-my-chat/src/models"
	"strings"
)

type Rule struct {
	Prompts   []string
	Action    ActionFunc
	NextState SessionState
}

type Session struct {
	UserID string
	State  SessionState
	Memory map[string]string
}

type ActionInput struct {
	Session *Session
	Message *models.Message
	Output  chan<- models.Message
}

type ActionFunc func(ctx context.Context, input ActionInput)
type MatcherFunc func(rules []Rule, input string) (Rule, bool)
type RuleEngineOption func(engine *RuleEngine)

type SessionRepo interface {
	GetOrCreate(ctx context.Context, sessionID string) *Session
}

func WithMatcher(m MatcherFunc) RuleEngineOption {
	return func(engine *RuleEngine) {
		engine.matcher = m
	}
}

func WithSessionRepo(repo SessionRepo) RuleEngineOption {
	return func(engine *RuleEngine) {
		engine.sessionRepo = repo
	}
}

type RuleEngine struct {
	matcher     MatcherFunc
	sessionRepo SessionRepo
	rules       []Rule
	ruleGroups  map[string][]Rule
}

func NewRuleEngine(opts ...RuleEngineOption) *RuleEngine {
	engine := &RuleEngine{}

	for _, opt := range opts {
		opt(engine)
	}

	if engine.matcher == nil {
		engine.matcher = DefaultMatcher
	}

	if engine.sessionRepo == nil {
		engine.sessionRepo = NewInMemorySessionRepo()
	}

	return engine
}

func (e *RuleEngine) RegisterRule(rule ...Rule) {
	e.rules = append(e.rules, rule...)
}

func (e *RuleEngine) HandleMessage(ctx context.Context, msg *models.Message, msgCh chan<- models.Message) {
	session := e.sessionRepo.GetOrCreate(ctx, msg.ID)
	actionInput := ActionInput{Session: session, Message: msg, Output: msgCh}

	switch state := session.State.(type) {
	case IdleState:
		e.handleIdleState(ctx, actionInput)
	case WaitingInputState:
		e.handleWaitingInputState(ctx, actionInput, state)
	case WaitingChoiceState:
		e.handleWaitingChoiceState(ctx, actionInput, state)
	default:
		e.handleUnknownState(actionInput)
	}
}

func (e *RuleEngine) handleIdleState(ctx context.Context, input ActionInput) {
	rule, ok := e.matcher(e.rules, input.Message.Input)
	if !ok {
		input.Message.Output = "desculpe não entendi"
		input.Output <- *input.Message
		return
	}

	// TODO manage session memory
	input.Session.State = rule.NextState
	rule.Action(ctx, input)

}

func (e *RuleEngine) handleWaitingInputState(ctx context.Context, input ActionInput, state WaitingInputState) {
	if strings.TrimSpace(input.Message.Input) == "" {
		input.Message.Output = "Por favor, responda ao que foi solicitado"
		input.Output <- *input.Message
		return
	}

	input.Session.State = IdleState{}
	state.Action(ctx, input)
}

func (e *RuleEngine) handleWaitingChoiceState(ctx context.Context, input ActionInput, state WaitingChoiceState) {
	handler, ok := state.Choices[input.Message.Input]
	if !ok {
		input.Message.Output = "Opção inválida. " + state.Prompt
		input.Output <- *input.Message
		return
	}

	input.Session.State = IdleState{}
	handler(ctx, input)

}

func (e *RuleEngine) handleUnknownState(input ActionInput) {
	input.Message.Output = "Erro interno: estado desconhecido."
	input.Output <- *input.Message
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

// func (e *RuleEngine) HandleMessage(msg models.Message, msgCh chan<- models.Message) {
// 	response := &msg
// 	actionPair := e.handleMsg(response, msgCh)
// 	e.checkOption = response.ResponseType == models.OptionResponse

// 	msgCh <- *response

// 	if actionPair != nil {
// 		e.actionSvc.Enqueue(*actionPair)
// 	}

// }

// func (e *RuleEngine) handleMsg(msg *models.Message, msgCh chan<- models.Message) *core.ActionReplyPair {
// 	if e.checkOption {
// 		return e.handleOption(msg, msgCh)
// 	}
// 	return e.handleIntent(msg, msgCh)
// }

// func (e *RuleEngine) handleIntent(response *models.Message, msgCh chan<- models.Message) *core.ActionReplyPair {
// 	intent, ok := e.intents.GetIntent(response.Input)

// 	if !ok {
// 		response.Output = core.MessageNoIntent
// 		response.ResponseType = models.TextResponse
// 		return nil
// 	}

// 	if intent.Action != nil {
// 		response.Output = core.MessageActionOngoing
// 		response.ResponseType = models.TextResponse
// 		return &core.ActionReplyPair{ReplyTo: msgCh, Action: *intent.Action, Input: *response}
// 	}

// 	if !intent.HasOptions() {
// 		response.ResponseType = models.TextResponse
// 		response.Output = intent.Response
// 		return nil
// 	}

// 	if e.cachedOptions.intent != intent.Key {
// 		NewCachedOptions(intent)
// 	}

// 	response.ResponseType = models.OptionResponse
// 	response.Options = utils.Map(intent.Options.items, func(o Option) models.Option {
// 		return models.Option{ID: o.Key, Name: o.Name}
// 	})

// 	return nil

// }

// func (e *RuleEngine) handleOption(msg *models.Message, msgCh chan<- models.Message) *core.ActionReplyPair {
// 	if opt, found := e.cachedOptions.GetOption(msg.Input); found {
// 		if opt.Action != nil {
// 			msg.Output = core.MessageActionOngoing
// 			msg.ResponseType = models.TextResponse
// 			return &core.ActionReplyPair{ReplyTo: msgCh, Action: *opt.Action, Input: *msg}
// 		}

// 		msg.Output = opt.Response
// 		msg.ResponseType = models.TextResponse
// 	}
// 	return nil
// }
