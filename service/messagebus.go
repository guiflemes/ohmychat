package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Handler interface {
	Handle(message Message) error
}

type handleMessage struct {
	Topic   string
	Handler Handler
}

type Handlers []handleMessage

func (h Handlers) GetHandler(topic string) func(message Message) error {
	for _, item := range h {
		if item.Topic == topic {
			return item.Handler.Handle
		}
	}
	return nil
}

func (h *Handlers) AddHandler(handler handleMessage) {
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
	wg       sync.WaitGroup
}

func NewBus() *messageBus {
	return &messageBus{
		queue: make(chan Message, 0),
		done:  make(chan struct{}),
	}
}

func (m *messageBus) SetHandler(topic string, handler Handler) {
	m.handlers = append(m.handlers, handleMessage{Topic: topic, Handler: handler})
}

func (m *messageBus) Consume() {

	m.wg.Add(1)
	go func() {
		m.wg.Done()
		for {
			select {

			case msg, ok := <-m.queue:
				if ok {
					m.handler(msg)
				}

			case <-m.done:
				return

			default:
			}
		}

	}()

}

func (m *messageBus) Publish(message Message) {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		select {
		case <-m.done:
			return
		case m.queue <- message:

		}
	}()
}

func (m *messageBus) Close() {
	close(m.done)
	m.wg.Wait()
	close(m.queue)
}

func (m *messageBus) handler(message Message) {

	handler := m.handlers.GetHandler(message.Meta().Topic)
	if handler == nil {
		log.Printf("handler to the given topic '%s' not found\n", message.Meta().Topic)
	}

	m.wg.Add(1)
	go func() {
		m.wg.Done()
		if err := handler(message); err != nil {
			log.Println(fmt.Errorf("Error '%s' handling message '%s' at topic '%s'", err, message.Meta().Id, message.Meta().Topic))
		}
	}()
}
