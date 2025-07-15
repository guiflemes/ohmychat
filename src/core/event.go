package core

import (
	"github.com/guiflemes/ohmychat/src/message"
	"time"
)

type EventType uint8

const (
	EventSuccess EventType = iota
	EventError
)

type Event struct {
	Type  EventType
	Msg   *message.Message
	Error error
	Time  time.Time
}

func (e *Event) WithError(err error) {
	e.Error = err
	e.Type = EventError
}

type OnEvent func(event Event)

func NewEvent(msg message.Message) *Event {
	return &Event{
		Msg:  &msg,
		Time: time.Now(),
	}
}

func NewEventError(err error) Event {
	return Event{
		Type:  EventError,
		Error: err,
		Time:  time.Now(),
	}
}

func NewEventErrorWithMessage(msg message.Message, err error) Event {
	return Event{
		Type:  EventError,
		Msg:   &msg,
		Error: err,
		Time:  time.Now(),
	}
}

func NewEventSuccess(msg message.Message) Event {
	return Event{
		Type:  EventError,
		Msg:   &msg,
		Error: nil,
		Time:  time.Now(),
	}
}

type EventHandlerOption func(h *EventHandler)

func EventWithMaxPool(maxPool uint8) EventHandlerOption {
	return func(h *EventHandler) {
		h.maxPool = maxPool
	}
}

func EventWithCallback(callback OnEvent) EventHandlerOption {
	return func(h *EventHandler) {
		h.onEvent = callback
	}
}

type EventHandler struct {
	maxPool uint8
	onEvent OnEvent
}

func NewEventHandler(options ...EventHandlerOption) *EventHandler {
	handler := &EventHandler{maxPool: 5}

	for _, opt := range options {
		opt(handler)
	}

	return handler
}

func (e *EventHandler) SetCallback(cb func(Event)) {
	e.onEvent = cb
}

func (h *EventHandler) Handler(cCtx *ChatContext, eventCh <-chan Event) {
	sem := make(chan struct{}, h.maxPool)
	for {
		select {
		case e := <-eventCh:
			go func(event Event) {
				sem <- struct{}{}
				defer func() { <-sem }()

				if h.onEvent != nil {
					h.onEvent(event)
				}
			}(e)
		case <-cCtx.Done():
			return
		}
	}
}
