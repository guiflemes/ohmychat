package service

import (
	"log"

	"github.com/google/uuid"
)

type Handler struct {
	Topic       string
	HandlerFunc func(message Message)
}

type Handlers []Handler

func (h Handlers) GetHandler(topic string) func(message Message) {
	for _, item := range h {
		if item.Topic == topic {
			return item.HandlerFunc
		}
	}
	return nil
}

func (h *Handlers) AddHandler(handler Handler) {
	for _, item := range *h {
		if item.Topic == handler.Topic {
			log.Printf(`Handler with the given topic %s already exists.
			Skipping new registration. The existing handler will be maintained.`, handler.Topic)
		}
		return
	}

	*h = append(*h, handler)
}

type MessageMeta struct {
	Id    uuid.UUID
	Topic string
}

type Message interface {
	Meta() MessageMeta
}

type messageBus struct {
	queue    chan Message
	handlers Handlers
	done     chan struct{}
}

func NewBus() *messageBus {
	return &messageBus{
		queue: make(chan Message, 0),
		done:  make(chan struct{}),
	}
}

func (m *messageBus) SetHandler(handler Handler) {
	m.handlers = append(m.handlers, handler)
}

func (m *messageBus) Consume() {
	for {
		select {

		case msg, ok := <-m.queue:
			if ok {
				m.handler(msg)
			}

		case <-m.done:
			close(m.queue)
			break

		default:
		}
	}

}

func (m *messageBus) Publish(message Message) {
	m.queue <- message
}

func (m *messageBus) Close() {
	close(m.done)
}

func (m *messageBus) handler(message Message) {

	handler := m.handlers.GetHandler(message.Meta().Topic)
	if handler == nil {
		log.Printf("handler to the given topic '%s' not found\n", message.Meta().Topic)
	}

	go handler(message)

}
