package context

import (
	"context"
	"errors"
	"oh-my-chat/src/message"
	"oh-my-chat/src/session"
)

type SessionAdapter interface {
	GetOrCreate(ctx context.Context, sessionID string) *session.Session
	Save(ctx context.Context, session *session.Session) error
}

type ChatContext struct {
	ctx            context.Context
	cancel         context.CancelFunc
	workflow       string
	metadata       map[string]any
	receiveCh      chan string
	outputCh       chan message.Message
	shutdownCh     chan struct{}
	session        *session.Session
	sessionAdapter SessionAdapter
}

func NewChatContext() *ChatContext {
	ctx, cancel := context.WithCancel(context.Background())

	return &ChatContext{
		ctx:            ctx,
		cancel:         cancel,
		metadata:       make(map[string]any),
		receiveCh:      make(chan string, 10),
		outputCh:       make(chan message.Message, 10),
		shutdownCh:     make(chan struct{}),
		sessionAdapter: session.NewInMemorySessionRepo(),
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

func (c *ChatContext) GetOrCreateSession(userID string) *session.Session {
	if c.session == nil {
		c.session = c.sessionAdapter.GetOrCreate(c.ctx, userID)
	}
	return c.session
}

func (c *ChatContext) SaveSession() error {
	if c.session == nil {
		return errors.New("no session to save")
	}
	return c.sessionAdapter.Save(c.ctx, c.session)
}

func (c *ChatContext) Session() *session.Session {
	return c.session
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
