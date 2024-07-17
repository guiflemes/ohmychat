package ohmychat

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"oh-my-chat/src/config"
	"oh-my-chat/src/connector/cli"
	"oh-my-chat/src/models"
)

func RunCli() {
	// use it to simulate up to cli be ready to run on app

	config := config.OhMyChatConfig{}
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-sigCh
		cancel()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bot := models.NewBot(config)
		conn, _ := cli.NewCliConnector(bot)
		conn.Acquire(ctx, make(chan<- models.Message))
	}()

	wg.Wait()
}
