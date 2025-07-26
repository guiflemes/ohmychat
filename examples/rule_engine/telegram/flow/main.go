package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	)

	tBot, err := tgbotapi.NewBotAPI("7879451742:AAExEruQ-EGx62fer25IFPYMdKV_Qu2OaBQ")
	if err != nil {
		log.Panicf("error starting telegram bot %s", err.Error())
	}
	chatBot := ohmychat.NewOhMyChat(connector.Telegram(tBot), ohmychat.WithEventCallback(logOnEvent))
	log.Println("running telegram bot...")
	chatBot.Run(e)
	log.Println("telegram bot finished")
}

func logOnEvent(event ohmychat.Event) {
	switch event.Type {
	case ohmychat.EventError:
		if event.Msg != nil {
			log.Printf("error on ohmychat '%s': %s", event.Msg.ID, event.Error.Error())
			return
		}
		log.Printf("error: %s", event.Error.Error())
	case ohmychat.EventSuccess:
		log.Printf("success on ohmychat '%s'", event.Msg.ID)
	}
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
		ThenSayAndContinue("só mais um pouco").
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
