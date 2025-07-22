package main

import (
	"fmt"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/engine/rule_engine"

	"regexp"

	"github.com/guiflemes/ohmychat/connector"
)

func main() {
	engine := rule_engine.NewRuleEngine()
	engine.RegisterRule(
		rule_engine.Rule{
			Prompts: []string{"fazer pedido", "enviar pedido"},
			Action: func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				msg.Output = "Qual o número do pedido?"
				ctx.SendOutput(msg)
			},
			NextState: ohmychat.WaitingInputState{
				PromptEmptyMessage: "Por favor, informe o número do pedido.",
				PromptExit:         "solicitação de pedido cancelado",
				ExitInput:          "sair",
				Action: ohmychat.WithValidation(
					func(input string) bool {
						match, _ := regexp.MatchString(`^PD:\s?\d{9}$`, input)
						return match
					},
					"Número de pedido inválido. Use o formato PD:123456789",
					func(ctx *ohmychat.Context, msg *ohmychat.Message) {
						msg.Output = fmt.Sprintf("Pedido %q registrado com sucesso!", msg.Input)
						ctx.SendOutput(msg)
					},
				),
			},
		},

		rule_engine.Rule{
			Prompts: []string{"ola", "ola tudo bem", "hello", "hello"},
			Action: func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				msg.Output = "Ola tudo bem e vc?"
				ctx.SendOutput(msg)
			},
			NextState: ohmychat.IdleState{},
		},

		rule_engine.Rule{
			Prompts: []string{"quero um cao", "cachorro", "dog"},
			Action: func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				msg.ResponseType = ohmychat.OptionResponse
				msg.Options = []ohmychat.Option{
					{ID: "beagle", Name: "beagle"},
					{ID: "pinscher", Name: "pinscher"},
					{ID: "pastor", Name: "pastor"},
					{ID: "pitbull", Name: "pitbull"},
				}
				ctx.SendOutput(msg)
			},
			NextState: ohmychat.WaitingChoiceState{
				Choices: ohmychat.Choices{
					"beagle": func(ctx *ohmychat.Context, msg *ohmychat.Message) {
						msg.Output = "legal, o cão mais fofo e gordo que existe"
						ctx.SendOutput(msg)
					},
					"pinscher": func(ctx *ohmychat.Context, msg *ohmychat.Message) {
						msg.Output = "legal, o cão mais feroz do mundo"
						ctx.SendOutput(msg)
					},
				}.BindMany(func(ctx *ohmychat.Context, msg *ohmychat.Message) {
					msg.Output = fmt.Sprintf("nossa seu cão %s é tao sem graça", msg.Input)
					ctx.SendOutput(msg)
				}, "pastor", "pitbull"),
			},
		},
	)

	chatBot := ohmychat.NewOhMyChat(connector.Cli())
	chatBot.Run(engine)
}
