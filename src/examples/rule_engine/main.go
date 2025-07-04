package main

import (
	"context"
	"fmt"
	"oh-my-chat/src/app"
	"oh-my-chat/src/core/rule_engine"
	"oh-my-chat/src/models"
	"regexp"
)

func main() {
	engine := rule_engine.NewRuleEngine()
	engine.RegisterRule(
		rule_engine.Rule{
			Prompts: []string{"fazer pedido", "enviar pedido"},
			Action: func(ctx context.Context, input rule_engine.ActionInput) {
				input.Message.Output = "Qual o número do pedido?"
				input.Output <- *input.Message
			},
			NextState: rule_engine.WaitingInputState{
				PromptEmptyMessage: "Por favor, informe o número do pedido.",
				Action: rule_engine.WithValidation(
					func(input string) bool {
						match, _ := regexp.MatchString(`^PD:\s?\d{9}$`, input)
						return match
					},
					"Número de pedido inválido. Use o formato PD:123456789",
					func(ctx context.Context, input rule_engine.ActionInput) {
						input.Message.Output = fmt.Sprintf("Pedido %q registrado com sucesso!", input.Message.Input)
						input.Output <- *input.Message
					},
				),
			},
		},

		rule_engine.Rule{
			Prompts: []string{"ola", "ola tudo bem", "hello", "hello"},
			Action: func(ctx context.Context, input rule_engine.ActionInput) {
				input.Message.Output = "Ola tudo bem e vc?"
				input.Output <- *input.Message
			},
			NextState: rule_engine.IdleState{},
		},

		rule_engine.Rule{
			Prompts: []string{"quero um cao", "cachorro", "dog"},
			Action: func(ctx context.Context, input rule_engine.ActionInput) {
				input.Message.ResponseType = models.OptionResponse
				input.Message.Options = []models.Option{{ID: "beagle", Name: "beagle"}, {ID: "pinscher", Name: "pinscher"}}
				input.Output <- *input.Message
			},
			NextState: rule_engine.WaitingChoiceState{
				Choices: rule_engine.Choices{
					"beagle": func(ctx context.Context, input rule_engine.ActionInput) {
						input.Message.Output = "legal, o cão mais fofo e gordo que existe"
						input.Output <- *input.Message
					},
					"pinscher": func(ctx context.Context, input rule_engine.ActionInput) {
						input.Message.Output = "legal, o cão mais feroz do mundo"
						input.Output <- *input.Message
					},
				},
			},
		},
	)

	app.RunV2(engine)

}
