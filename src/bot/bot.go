package bot

import (
	"oh-my-chat/src/connector"
	"oh-my-chat/src/context"
	"oh-my-chat/src/core"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/message"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
)

type Bot struct {
	Connector connector.Connector
}

func (b *Bot) Run(engine core.Engine) {

	inputMsg := make(chan message.Message, 1)
	outputMsg := make(chan message.Message, 1)

	chatCtx := context.NewChatContext()

	processor := core.NewProcessor(engine)
	connector := core.NewMuitiChannelConnector(b.Connector)

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)

	//TODO  when cli connector is running, the cancelation never comes here, find a way out to fix it
	go func() {
		sig := <-sign
		logger.Logger.Info("Received signal, stopping", zap.String("signal", sig.String()))
		chatCtx.Shutdown()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		processor.Process(chatCtx.Context(), inputMsg, outputMsg)
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
