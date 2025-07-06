package context

import (
	"context"
	"oh-my-chat/src/message"
)

type ChatContext struct {
	ctx        context.Context
	cancel     context.CancelFunc
	workflow   string
	metadata   map[string]any
	receiveCh  chan string
	outputCh   chan message.Message
	shutdownCh chan struct{}
}

func NewChatContext() *ChatContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &ChatContext{
		ctx:        ctx,
		cancel:     cancel,
		metadata:   make(map[string]any),
		receiveCh:  make(chan string, 10),
		outputCh:   make(chan message.Message, 10),
		shutdownCh: make(chan struct{}),
	}
}

func (c *ChatContext) Context() context.Context {
	return c.ctx
}

func (c *ChatContext) Shutdown() {
	c.cancel()
	close(c.shutdownCh)
}

func (c *ChatContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *ChatContext) IsActive() bool {
	select {
	case <-c.ctx.Done():
		return false
	default:
		return true
	}
}

func (c *ChatContext) Set(key string, value any) {
	c.metadata[key] = value
}

func (c *ChatContext) Get(key string) (any, bool) {
	v, ok := c.metadata[key]
	return v, ok
}

// TODO use it to make api easier
// func (c *ChatContext) ReceiveChannel() <-chan string {
// 	return c.receiveCh
// }

// func (c *ChatContext) OutputChannel() chan<- message.Message {
// 	return c.outputCh
// }

// func (c *ChatContext) SendInput(input string) {
// 	c.receiveCh <- input
// }

// func (c *ChatContext) SendOutput(msg message.Message) {
// 	c.outputCh <- msg
// }
