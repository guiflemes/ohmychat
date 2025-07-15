package main

import (
	"fmt"
	"github.com/guiflemes/ohmychat/src/bot"
	"github.com/guiflemes/ohmychat/src/core/rule_engine"

	"github.com/guiflemes/ohmychat/src/core"
	"github.com/guiflemes/ohmychat/src/message"

	"regexp"

	"github.com/guiflemes/ohmychat/src/connector/cli"
)

func main() {
	engine := rule_engine.NewRuleEngine()
	engine.RegisterRule(
		rule_engine.Rule{
			Prompts: []string{"fazer pedido", "enviar pedido"},
			Action: func(ctx *core.Context, msg *message.Message) {
				msg.Output = "Qual o número do pedido?"
				ctx.SendOutput(msg)
			},
			NextState: core.WaitingInputState{
				PromptEmptyMessage: "Por favor, informe o número do pedido.",
				PromptExit:         "solicitação de pedido cancelado",
				ExitInput:          "sair",
				Action: core.WithValidation(
					func(input string) bool {
						match, _ := regexp.MatchString(`^PD:\s?\d{9}$`, input)
						return match
					},
					"Número de pedido inválido. Use o formato PD:123456789",
					func(ctx *core.Context, msg *message.Message) {
						msg.Output = fmt.Sprintf("Pedido %q registrado com sucesso!", msg.Input)
						ctx.SendOutput(msg)
					},
				),
			},
		},

		rule_engine.Rule{
			Prompts: []string{"ola", "ola tudo bem", "hello", "hello"},
			Action: func(ctx *core.Context, msg *message.Message) {
				msg.Output = "Ola tudo bem e vc?"
				ctx.SendOutput(msg)
			},
			NextState: core.IdleState{},
		},

		rule_engine.Rule{
			Prompts: []string{"quero um cao", "cachorro", "dog"},
			Action: func(ctx *core.Context, msg *message.Message) {
				msg.ResponseType = message.OptionResponse
				msg.Options = []message.Option{
					{ID: "beagle", Name: "beagle"},
					{ID: "pinscher", Name: "pinscher"},
					{ID: "pastor", Name: "pastor"},
					{ID: "pitbull", Name: "pitbull"},
				}
				ctx.SendOutput(msg)
			},
			NextState: core.WaitingChoiceState{
				Choices: core.Choices{
					"beagle": func(ctx *core.Context, msg *message.Message) {
						msg.Output = "legal, o cão mais fofo e gordo que existe"
						ctx.SendOutput(msg)
					},
					"pinscher": func(ctx *core.Context, msg *message.Message) {
						msg.Output = "legal, o cão mais feroz do mundo"
						ctx.SendOutput(msg)
					},
				}.BindMany(func(ctx *core.Context, msg *message.Message) {
					msg.Output = fmt.Sprintf("nossa seu cão %s é tao sem graça", msg.Input)
					ctx.SendOutput(msg)
				}, "pastor", "pitbull"),
			},
		},
	)

	chatBot := bot.NewOhMyChat(cli.NewCliConnector())
	chatBot.Run(engine)
}
