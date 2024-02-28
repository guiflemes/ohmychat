package service

import (
	"fmt"
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

type Message struct {
	Type MessageType
}

type messageBus struct {
	queue chan Message
	done  <-chan struct{}
}

func NewBus(queue chan Message, done chan struct{}) *messageBus {
	return &messageBus{queue: queue, done: done}
}

func (m *messageBus) Handler() {

	for {
		time.Sleep(time.Second)
		fmt.Println("Waiting a message")
		select {
		case <-m.done:
			select {
			case msg := <-m.queue:
				m.handler(msg)
			default:
				close(m.queue)
				return
			}
		case msg := <-m.queue:
			m.handler(msg)

		default:
		}

	}

}

func (m *messageBus) handler(message Message) {

	switch message.Type {

	case Event:
		fmt.Println("event")
	case Command:
		fmt.Println("command")

	}
}

func RunBus() {
	queue := make(chan Message, 0)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{}, 1)

	bus := NewBus(queue, done)
	go bus.Handler()

	limit := 10

	var wg sync.WaitGroup

	for i := 0; i < limit; i++ {
		message := Message{Type: Event}

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
		<-signalCh
		fmt.Println("Stoping gracefully")
		close(done)
	}()

	wg.Wait()

}
