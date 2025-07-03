package main

import (
	"context"
	"fmt"
	"oh-my-chat/src/app"
	"oh-my-chat/src/core/rule_engine"
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
				Prompt: "Por favor, informe o número do pedido.",
				Action: func(ctx context.Context, input rule_engine.ActionInput) {
					input.Message.Output = fmt.Sprintf("Pedido %q registrado com sucesso!", input.Message.Input)
					input.Output <- *input.Message
				},
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
				input.Message.Output = "vc quer um cão chamado marvao?"
				input.Output <- *input.Message
			},
			NextState: rule_engine.IdleState{},
		},
	)

	app.RunV2(engine)

}
