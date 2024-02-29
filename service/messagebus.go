package service

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type MessageType int

const (
	Event   MessageType = iota
	Command             = iota
)

type EventHandler struct {
	Topic   string
	Handler func()
}

type EventHandlers []EventHandler

func (h EventHandlers) GetHandler(topic string) func() {
	for _, e := range h {
		if e.Topic == topic {
			return e.Handler
		}
	}
	return nil
}

type MessageMeta struct {
	Topic string
	Type  MessageType
}

type Message interface {
	Meta() MessageMeta
}

type messageBus struct {
	queue         chan Message
	done          <-chan struct{}
	eventHandlers EventHandlers
}

func NewBus(queue chan Message, done chan struct{}) *messageBus {
	return &messageBus{
		queue: queue,
		done:  done}
}

func (m *messageBus) SetEventHandler(handler EventHandler) {
	m.eventHandlers = append(m.eventHandlers, handler)
}

func (m *messageBus) Handler() {

	for {
		time.Sleep(time.Second)
		fmt.Println("Waiting a message")
		select {
		case <-m.done:
			close(m.queue)
			return
		case msg := <-m.queue:
			m.handler(msg)

		default:
		}

	}

}

func (m *messageBus) handler(message Message) {

	switch message.Meta().Type {

	case Event:
		m.handlerEvent(message)
	case Command:
		fmt.Println("command handler not implemented yet")

	}
}

func (m *messageBus) handlerEvent(message Message) {
	handler := m.eventHandlers.GetHandler(message.Meta().Topic)

	if handler == nil {
		log.Printf("handler to the given topic '%s' not found\n", message.Meta().Topic)
	}

	go handler()
}

type SomeMessage struct {
	Type MessageType
}

func (s *SomeMessage) Meta() MessageMeta {
	return MessageMeta{
		Type:  s.Type,
		Topic: "some_name",
	}
}

func RunBus() {
	queue := make(chan Message, 0)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{}, 1)

	bus := NewBus(queue, done)
	bus.SetEventHandler(EventHandler{Topic: "some_name", Handler: func() { fmt.Println("handler event") }})

	go bus.Handler()

	limit := 10

	var wg sync.WaitGroup

	for i := 0; i < limit; i++ {
		message := &SomeMessage{}

		if i%2 == 0 {
			message.Type = Command
		}

		wg.Add(1)
		go func(msg Message) {
			queue <- msg
			wg.Done()
		}(message)
	}

	go func() {
		wg.Wait()

	}()

	select {
	case <-signalCh:
		close(done)
		return
	}
}
