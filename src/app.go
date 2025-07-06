package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"oh-my-chat/src/bot"
	settings "oh-my-chat/src/config"
	"oh-my-chat/src/connector/cli"
	"oh-my-chat/src/core"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/message"
)

func Run(engine core.Engine) {
	cfg := settings.OhMyChatConfig{
		Connector: settings.Connector{Provider: settings.Cli},
	}

	inputMsg := make(chan message.Message, 1)
	outputMsg := make(chan message.Message, 1)

	bot := bot.NewBot(cfg)
	ctx := bot.Ctx()

	processor := core.NewProcessor(engine)
	connector := core.NewMuitiChannelConnector(cli.NewCliConnector(bot))

	sign := make(chan os.Signal, 1)
	signal.Notify(sign, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)

	//TODO  when cli connector is running, the cancelation never comes here, find a way out to fix it
	go func() {
		sig := <-sign
		logger.Logger.Info("Received signal, stopping", zap.String("signal", sig.String()))
		bot.Shutdown()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		processor.Process(ctx, inputMsg, outputMsg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		connector.Request(ctx, inputMsg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		connector.Response(ctx, outputMsg)
	}()

	wg.Wait()

}
