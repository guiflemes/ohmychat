package core

import (
	"context"
	"oh-my-chat/src/message"
	"time"
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
	workflow       string
	metadata       map[string]any
	shutdownCh     chan struct{}
	sessionAdapter SessionAdapter
}

func NewChatContext(options ...ChatContextOption) *ChatContext {
	ctx, cancel := context.WithCancel(context.Background())

	chatCtx := &ChatContext{
		ctx:        ctx,
		cancel:     cancel,
		metadata:   make(map[string]any),
		shutdownCh: make(chan struct{}),
	}

	for _, opt := range options {
		opt(chatCtx)
	}

	if chatCtx.sessionAdapter == nil {
		chatCtx.sessionAdapter = NewInMemorySessionRepo()
	}

	return chatCtx
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

func (c *ChatContext) NewChildContext(msg message.Message, outputCh chan<- message.Message) (*Context, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 60*time.Second)

	sess, err := c.sessionAdapter.GetOrCreate(ctx, msg.User.ID)
	if err != nil {
		cancel()
		return nil, err
	}

	return &Context{
		ctx:      ctx,
		cancel:   cancel,
		parent:   c,
		session:  sess,
		outputCh: outputCh,
	}, nil
}

type Context struct {
	ctx      context.Context
	cancel   context.CancelFunc
	session  *Session
	parent   *ChatContext
	outputCh chan<- message.Message
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

func (c *Context) SendOutput(msg *message.Message) {
	c.parent.sessionAdapter.Save(c.Context(), c.session)
	c.outputCh <- *msg
}
