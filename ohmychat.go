package ohmychat

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Engine interface {
	HandleMessage(*Context, *Message)
}

type ohMyChat struct {
	connector    Connector
	eventHandler *EventHandler
}

type OhMyChatOption func(*ohMyChat)

func WithEventCallback(cb func(Event)) OhMyChatOption {
	return func(b *ohMyChat) {
		b.eventHandler.SetCallback(cb)
	}
}

func NewOhMyChat(connector Connector, opts ...OhMyChatOption) *ohMyChat {
	b := &ohMyChat{
		connector:    connector,
		eventHandler: NewEventHandler(),
	}

	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *ohMyChat) Run(engine Engine) {

	inputMsg := make(chan Message, 10)
	outputMsg := make(chan Message, 10)
	eventCh := make(chan Event, 10)

	chatCtx := NewChatContext(eventCh)
	chatCtx.InputCh = inputMsg
	processor := NewProcessor(engine)
	connector := NewMuitiChannelConnector(b.connector)

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
