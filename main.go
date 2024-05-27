package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"oh-my-chat/src/config"
	"oh-my-chat/src/core"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
	"oh-my-chat/src/services"
	"oh-my-chat/src/storage"
	"oh-my-chat/src/storage/action"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Run()
}

func Run() {

	var (
		inputMsg  = make(chan models.Message, 1)
		outputMsg = make(chan models.Message, 1)
	)
	ctx, cancel := context.WithCancel(context.Background())
	bot := models.NewBot(models.Telegram)

	storageService := &action.StorageActionMessage{}

	guidedEngine := core.NewGuidedResponseEngine(storageService, storage.NewLoadFileRepository())

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

	wg.Add(4)

	go func() {
		defer wg.Done()
		services.RunWorker(ctx, config.Worker{Number: 1}, storageService)
	}()

	go func() {
		defer wg.Done()
		processor.Process(ctx, inputMsg, outputMsg)
	}()

	go func() {
		defer wg.Done()
		connector.Request(ctx, inputMsg)
	}()

	go func() {
		defer wg.Done()
		connector.Response(ctx, outputMsg)
	}()

	wg.Wait()
}
