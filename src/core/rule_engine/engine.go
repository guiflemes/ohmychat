package rule_engine

import (
	"oh-my-chat/src/core"
	"oh-my-chat/src/models"
	"oh-my-chat/src/utils"
)

type RuleEngine struct {
	actionSvc     core.ActionStorageService
	intents       Intents
	cachedOptions CachedOptions
	checkOption   bool
}

func (e *RuleEngine) Name() string                                       { return "rule_engine" }
func (e *RuleEngine) GetActionStorageService() core.ActionStorageService { return e.actionSvc }

func (e *RuleEngine) HandleMessage(msg models.Message, msgCh chan<- models.Message) {
	response := &msg
	actionPair := e.handleMsg(response, msgCh)
	e.checkOption = response.ResponseType == models.OptionResponse

	msgCh <- *response

	if actionPair != nil {
		e.actionSvc.Enqueue(*actionPair)
	}

}

func (e *RuleEngine) handleMsg(msg *models.Message, msgCh chan<- models.Message) *core.ActionReplyPair {
	if e.checkOption {
		return e.handleOption(msg, msgCh)
	}
	return e.handleIntent(msg, msgCh)
}

func (e *RuleEngine) handleIntent(response *models.Message, msgCh chan<- models.Message) *core.ActionReplyPair {
	intent, ok := e.intents.GetIntent(response.Input)

	if !ok {
		response.Output = "desculpa não entendi"
		response.ResponseType = models.TextResponse
		return nil
	}

	if intent.Action != nil {
		response.Output = "processando action"
		response.ResponseType = models.TextResponse
		return &core.ActionReplyPair{ReplyTo: msgCh, Action: *intent.Action, Input: *response}
	}

	if !intent.HasOptions() {
		response.ResponseType = models.TextResponse
		response.Output = intent.Response
		return nil
	}

	if e.cachedOptions.intent != intent.Key {
		NewCachedOptions(intent)
	}

	response.ResponseType = models.OptionResponse
	response.Options = utils.Map(intent.Options.items, func(o Option) models.Option {
		return models.Option{ID: o.Key, Name: o.Name}
	})

	return nil

}

func (e *RuleEngine) handleOption(msg *models.Message, msgCh chan<- models.Message) *core.ActionReplyPair {
	if opt, found := e.cachedOptions.GetOption(msg.Input); found {
		if opt.Action != nil {
			msg.Output = "Processando action"
			msg.ResponseType = models.TextResponse
			return &core.ActionReplyPair{ReplyTo: msgCh, Action: *opt.Action, Input: *msg}
		}

		msg.Output = opt.Response
		msg.ResponseType = models.TextResponse
	}
	return nil
}
