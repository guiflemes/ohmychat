package service

import (
	"fmt"
	"log"
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
	queue         Queue
	eventHandlers EventHandlers
}

func NewBus(queue Queue) *messageBus {
	return &messageBus{
		queue: queue,
	}
}

func (m *messageBus) SetEventHandler(handler EventHandler) {
	m.eventHandlers = append(m.eventHandlers, handler)
}

func (m *messageBus) Handler() {

	for {
		msg, ok := m.queue.Consume()
		if ok {
			m.handler(msg)
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
