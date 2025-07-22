package main

import (
	"fmt"
	"strings"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/connector"
	"github.com/guiflemes/ohmychat/engine/rule_engine"
)

func main() {
	e := rule_engine.NewRuleEngine()
	e.RegisterRule(
		rule_engine.Rule{
			Prompts: []string{"camiseta", "comprar camiseta", "camisa"},
			Action: func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				msg.Output = "Olá! Vamos comprar uma camiseta. Qual marca você prefere? (Nike, Adidas)"
				ctx.SendOutput(msg)
			},
			NextState: startFlow(),
		},
	)

	chatBot := ohmychat.NewOhMyChat(connector.Cli())
	chatBot.Run(e)
}

func startFlow() ohmychat.WaitingInputState {
	return ohmychat.WaitingInputState{
		Action: ohmychat.WithValidation(func(input string) bool {
			ii := strings.ToLower(input)
			if ii != "nike" && ii != "adidas" {
				return false
			}
			return true
		},
			"nao temos essa marca",
			func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				ctx.Session().Memory["marca"] = strings.ToLower(msg.Input)
				msg.Output = "legal agora escolha seu tamanho"
				msg.ResponseType = ohmychat.OptionResponse
				msg.Options = []ohmychat.Option{
					{Name: "P", ID: "P"},
					{Name: "M", ID: "M"},
					{Name: "G", ID: "G"},
				}

				ctx.SetSessionState(ohmychat.WaitingChoiceState{
					PromptInvalidOption: "opção invalida",
					Prompt:              "Qual o tamanho você gostaria",
					Choices: ohmychat.Choices{}.BindMany(
						func(ctx *ohmychat.Context, msg *ohmychat.Message) {
							ctx.Session().Memory["size"] = msg.Input
							msg.Output = "legal vamos escolher uma cor agora"
							msg.Options = []ohmychat.Option{
								{Name: "blue", ID: "blue"},
								{Name: "red", ID: "red"},
								{Name: "write", ID: "write"},
							}
							ctx.SendOutput(msg)

							ctx.SetSessionState(ohmychat.WaitingChoiceState{
								PromptInvalidOption: "opçaõ invalida",
								Prompt:              "Qual a cor você gostaria?",
								Choices: ohmychat.Choices{}.BindMany(
									func(ctx *ohmychat.Context, msg *ohmychat.Message) {
										brand := ctx.Session().Memory["marca"]
										size := ctx.Session().Memory["size"]
										msg.Output = fmt.Sprintf("leval entao vc quer um camisa da %s, tamanho %s e cor %s", brand, size, msg.Input)
										ctx.SendOutput(msg)
									}, "red", "blue", "write"),
							})

						}, "P", "M", "G"),
				})

				ctx.SendOutput(msg)
			}),
	}
}
