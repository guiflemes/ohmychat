package bot

import (
	"github.com/guiflemes/ohmychat/src/core"
	"github.com/guiflemes/ohmychat/src/message"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type ohMyChat struct {
	Connector    core.Connector
	eventHandler *core.EventHandler
}

type OhMyChatOption func(*ohMyChat)

func WithEventCallback(cb func(core.Event)) OhMyChatOption {
	return func(b *ohMyChat) {
		b.eventHandler.SetCallback(cb)
	}
}

func NewOhMyChat(connector core.Connector, opts ...OhMyChatOption) *ohMyChat {
	b := &ohMyChat{
		Connector:    connector,
		eventHandler: core.NewEventHandler(),
	}

	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *ohMyChat) Run(engine core.Engine) {

	inputMsg := make(chan message.Message, 10)
	outputMsg := make(chan message.Message, 10)
	eventCh := make(chan core.Event, 10)

	chatCtx := core.NewChatContext(eventCh)
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

	wg.Add(1)
	go func() {
		defer wg.Done()
		b.eventHandler.Handler(chatCtx, eventCh)
	}()

	wg.Wait()

}
