package main

import (
	"fmt"
	"log"

	"github.com/guiflemes/ohmychat"
	"github.com/guiflemes/ohmychat/core/rule_engine"

	"github.com/guiflemes/ohmychat/core"
	"github.com/guiflemes/ohmychat/message"

	"regexp"

	"github.com/guiflemes/ohmychat/connector/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
				ExitInput:          "exit",
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
				msg.Output = "escolha sua opção de doguinho"
				msg.Options = []message.Option{{ID: "beagle", Name: "beagle"}, {ID: "pinscher", Name: "pinscher"}}
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
				},
			},
		},
	)

	tBot, err := tgbotapi.NewBotAPI("YOUR_TOKEN")
	if err != nil {
		log.Panicf("error starting telegram bot %s", err.Error())
	}
	chatBot := ohmychat.NewOhMyChat(telegram.NewTelegramConnector(tBot), ohmychat.WithEventCallback(logOnEvent))
	log.Println("running telegram bot...")
	chatBot.Run(engine)
	log.Println("telegram bot finished")
}

func logOnEvent(event core.Event) {
	switch event.Type {
	case core.EventError:
		if event.Msg != nil {
			log.Printf("error on message '%s': %s", event.Msg.ID, event.Error.Error())
			return
		}
		log.Printf("error: %s", event.Error.Error())
	case core.EventSuccess:
		log.Printf("success on message '%s'", event.Msg.ID)
	}
}
