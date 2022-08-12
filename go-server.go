package goserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	server *http.Server
	wg     *errgroup.Group
	once   sync.Once
	ctx    context.Context
	cancel context.CancelFunc
}

func NewServer(server *http.Server, ctx context.Context, cancel context.CancelFunc) *Server {
	return &Server{
		server: server,
		wg:     &errgroup.Group{},
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Server) Start(waitFuncs ...WGFunc) {
	s.once.Do(
		func() {
			s.start(waitFuncs...)
		})
}

func (s *Server) StartAndWait(waitFuncs ...WGFunc) error {
	s.once.Do(
		func() {
			s.start(waitFuncs...)
		})
	return s.wg.Wait()
}

func (s *Server) Wait() error {
	return s.wg.Wait()
}

func (s *Server) start(waitFuncs ...WGFunc) {
	for _, waitFunc := range waitFuncs {
		s.wg.Go(waitFunc)
	}
	s.wg.Go(func() error {
		defer func() {
			s.cancel()
		}()
		return s.server.ListenAndServe()
	})
}

type WGFunc func() error

// WatchContext watches the provided context. When the context is Done, it will begin to shutdown the server and give time for other processes to finish.
func (s *Server) WatchContext() WGFunc {
	return func() error {
		<-s.ctx.Done()
		return gracefulShutdown(s.server)
	}
}

func (s *Server) WatchSignal(signals ...os.Signal) WGFunc {
	return func() error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, signals...)
		return gracefulShutdown(s.server)
	}
}

func gracefulShutdown(server *http.Server) error {
	closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(closeCtx); err != nil {
		return err
	}
	return closeCtx.Err()
}
