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

	flow := buidlerFlow()
	e.RegisterRule(
		rule_engine.Rule{
			Prompts: []string{"camiseta", "comprar camiseta", "camisa"},
			Action:  flow.Start(),
		},

		rule_engine.Rule{
			Prompts: []string{"roupa", "pano", "vestuario"},
			Action: func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				msg.Output = "Olá! Vamos comprar uma camiseta. Qual marca você prefere? (Nike, Adidas)"
				ctx.SendOutput(msg)
			},
			NextState: chainingFlow(),
		},
	)

	chatBot := ohmychat.NewOhMyChat(connector.Cli())
	chatBot.Run(e)
}

func buidlerFlow() *rule_engine.FlowBuilder {
	return rule_engine.NewFlow().
		AskChoice(
			"Olá! Vamos comprar uma camiseta. Qual marca você prefere?",
			[]string{"Nike", "Adidas"},
			func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				ctx.Session().Memory["marca"] = strings.ToLower(msg.Input)
			},
		).
		AskChoice(
			"Qual tamanho voce gostaria",
			[]string{"P", "M", "G"},
			func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				ctx.Session().Memory["size"] = msg.Input
			},
		).
		AskChoice(
			"Qual a cor escolhida",
			[]string{"blue", "red", "write"},
			func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				ctx.Session().Memory["color"] = msg.Input
			},
		).
		ThenSayAndContinue("preparando tudo, só um momento...").
		ThenSayAndContinue("só mais um instante...").
		ThenFinal(
			func(ctx *ohmychat.Context, msg *ohmychat.Message) {
				brand := ctx.Session().Memory["marca"]
				size := ctx.Session().Memory["size"]
				color := ctx.Session().Memory["color"]
				msg.Output = fmt.Sprintf("leval entao vc quer um camisa da %s, tamanho %s e cor %s", brand, size, color)
				ctx.SendOutput(msg)
			},
		)
}

func chainingFlow() ohmychat.WaitingInputState {
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
										msg.BotMode = true // set bot mode to lock user reply
										ctx.SendOutput(msg)
										ctx.SetSessionState(ohmychat.WaitingBotResponseState{
											OnDone: func(ctx *ohmychat.Context, msg *ohmychat.Message) {
												msg.Output = "legal compra feita"
												ctx.SendOutput(msg)
												//set idle state to release user reply
												ctx.SetSessionState(ohmychat.IdleState{})
											},
										})
									}, "red", "blue", "write"),
							})

						}, "P", "M", "G"),
				})

				ctx.SendOutput(msg)
			}),
	}
}
