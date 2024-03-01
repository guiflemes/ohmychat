package service

import (
	"log"
)

type Queue struct {
	message chan Message
	done    chan struct{}
}

func (q *Queue) Publish(msg Message) {
	log.Printf("publishing message at topic %s", msg.Meta().Topic)
	q.message <- msg
}

func (q *Queue) Consume() (Message, bool) {
	for {
		select {
		case msg, ok := <-q.message:
			return msg, ok
		case <-q.done:
			close(q.message)
		}
	}
}
