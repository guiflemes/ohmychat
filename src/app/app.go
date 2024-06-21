package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"oh-my-chat/src/api"
	"oh-my-chat/src/config"
	"oh-my-chat/src/core"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
	"oh-my-chat/src/services"
	"oh-my-chat/src/storage"
	"oh-my-chat/src/storage/action"
)

func Run(config config.OhMyChatConfig) {
	// TODO inject it
	mux := api.NewHttpMux(api.NewOhMyChatApi())

	{
		var (
			inputMsg  = make(chan models.Message, 1)
			outputMsg = make(chan models.Message, 1)
		)
		ctx, cancel := context.WithCancel(context.Background())
		bot := models.NewBot(config.Connector)

		storageService := action.NewMemoryQueue()

		guidedEngine := core.NewGuidedResponseEngine(
			storageService,
			storage.NewLoadFileRepository(),
		)

		processor := core.NewProcessor(storage.NewMemoryChatbotRepo(), core.Engines{guidedEngine})
		connector := core.NewMuitiChannelConnector(bot)

		sign := make(chan os.Signal, 1)
		signal.Notify(sign, syscall.SIGTERM, os.Interrupt)

		go func() {
			sig := <-sign
			logger.Logger.Info("Received signal, stopping", zap.String("signal", sig.String()))
			cancel()
		}()

		var wg sync.WaitGroup

		if config.Api.Enabled {
			wg.Add(1)
			go func() {
				defer wg.Done()
				api.RunApi(ctx, config.Api, mux)
			}()
		}

		if config.Worker.Enabled {
			wg.Add(1)
			go func() {
				defer wg.Done()
				services.RunWorker(ctx, config.Worker, storageService)
			}()
		}

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
}
