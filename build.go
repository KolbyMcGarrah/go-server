package goserver

import (
	"context"
	"net/http"
)

type Builder struct {
	config  *Config
	server  *http.Server
	handler http.Handler
	options Options
	ctx     context.Context
	cancel  context.CancelFunc
}

func NeWBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithCustomServer(server *http.Server) *Builder {
	b.server = server
	return b
}

func (b *Builder) WithConfig(config *Config) *Builder {
	b.config = config
	return b
}

func (b *Builder) WithOptions(options ...Option) *Builder {
	b.options = append(b.options, options...)
	return b
}

func (b *Builder) WithContext(ctx context.Context, cancel context.CancelFunc) *Builder {
	b.cancel = cancel
	b.ctx = ctx
	return b
}

func (b *Builder) WithHandler(handler *http.Handler) *Builder {
	b.handler = *handler
	return b
}

func (b *Builder) Build(ctx context.Context) *Server {

	if b.server == nil {
		b.server = &http.Server{}
		b.server.Handler = b.handler
	}

	b.options.Apply(b.server)

	if b.ctx == nil || b.cancel == nil {
		b.ctx, b.cancel = context.WithCancel(context.Background())
	}

	return NewServer(b.server, b.ctx, b.cancel)
}
