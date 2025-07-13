package bot

import (
	"oh-my-chat/src/core"
	"oh-my-chat/src/message"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Bot struct {
	Connector core.Connector
}

func (b *Bot) Run(engine core.Engine) {

	inputMsg := make(chan message.Message, 1)
	outputMsg := make(chan message.Message, 1)

	chatCtx := core.NewChatContext()

	processor := core.NewProcessor(engine)
	connector := core.NewMuitiChannelConnector(b.Connector)

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)

	go func() {
		<-sign
		chatCtx.Shutdown()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		processor.Process(chatCtx, inputMsg, outputMsg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		connector.Request(chatCtx, inputMsg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		connector.Response(chatCtx, outputMsg)
	}()

	wg.Wait()

}
