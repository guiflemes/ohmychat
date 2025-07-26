//go:generate mockgen -source context.go -destination ./mocks/context.go -package mocks
package ohmychat

import (
	"context"
	"time"
)

const (
	ReplyDispatched = 1 << 0
)

type SessionAdapter interface {
	GetOrCreate(ctx context.Context, sessionID string) (*Session, error)
	Save(ctx context.Context, session *Session) error
}

type ChatContextOption func(ctx *ChatContext)

func WithSessionAdapter(adapater SessionAdapter) ChatContextOption {
	return func(ctx *ChatContext) {
		ctx.sessionAdapter = adapater
	}
}

type ChatContext struct {
	ctx            context.Context
	cancel         context.CancelFunc
	metadata       map[string]any
	shutdownCh     chan struct{}
	eventCh        chan<- Event
	sessionAdapter SessionAdapter
}

func NewChatContext(eventCh chan<- Event, options ...ChatContextOption) *ChatContext {
	ctx, cancel := context.WithCancel(context.Background())

	chatCtx := &ChatContext{
		ctx:        ctx,
		cancel:     cancel,
		metadata:   make(map[string]any),
		shutdownCh: make(chan struct{}),
		eventCh:    eventCh,
	}

	for _, opt := range options {
		opt(chatCtx)
	}

	if chatCtx.sessionAdapter == nil {
		chatCtx.sessionAdapter = NewInMemorySessionRepo()
	}

	return chatCtx
}

func (c *ChatContext) SendEvent(event Event) {
	c.eventCh <- event
}

func (c *ChatContext) SaveSession(ctx context.Context, session *Session) error {
	session.LastActivityAt = time.Now()
	err := c.sessionAdapter.Save(ctx, session)
	if err != nil {
		c.SendEvent(NewEventError(err))
	}
	return err
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

func (c *ChatContext) NewChildContext(msg Message, outputCh chan<- Message) (*Context, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 60*time.Second)

	sess, err := c.sessionAdapter.GetOrCreate(ctx, msg.User.ID)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Context{
		ctx:             ctx,
		cancel:          cancel,
		parent:          c,
		session:         sess,
		outputCh:        outputCh,
		replyDispatched: ReplyDispatched,
	}, nil
}

type Context struct {
	ctx             context.Context
	cancel          context.CancelFunc
	session         *Session
	parent          *ChatContext
	outputCh        chan<- Message
	replyDispatched uint8
}

func (c *Context) Context() context.Context {
	return c.ctx
}

func (c *Context) Cancel() {
	c.cancel()
}

func (c *Context) IsActive() bool {
	select {
	case <-c.ctx.Done():
		return false
	default:
		return true
	}
}

func (c *Context) Session() *Session {
	return c.session
}

func (c *Context) SetSessionState(state SessionState) {
	c.session.State = state
}

func (c *Context) MessageHasBeenReplyed() bool {
	return c.replyDispatched != 0
}

func (c *Context) SendOutput(msg *Message) {
	c.parent.SaveSession(c.Context(), c.session)
	c.replyDispatched |= ReplyDispatched
	c.outputCh <- *msg
}
