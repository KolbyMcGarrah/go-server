package goserver

import (
	"context"
	"log"
	"net"
	"net/http"
)

// SetTimeOuts sets the ReadTimeout, ReadHeaderTimeout, WriteTimeout, and IdleTimeout of the server from the provided config unset timeouts will default to 0 or No Timeout
func SetTimeOuts(config *Config) ApplyFunc {
	return func(s *http.Server) {
		s.ReadHeaderTimeout = config.ReadHeaderTimeout
		s.ReadTimeout = config.ReadTimeout
		s.WriteTimeout = config.WriteTimeout
		s.IdleTimeout = config.IdleTimeout
	}
}

// SetLogger sets the logger for the server
func SetLogger(logger *log.Logger) ApplyFunc {
	return func(s *http.Server) {
		s.ErrorLog = logger
	}
}

// SetBaseContext sets the basecontext for the server and currently ignores the net.Listener
func SetBaseContext(ctx context.Context) ApplyFunc {
	return func(s *http.Server) {
		s.BaseContext = func(_ net.Listener) context.Context {
			return ctx
		}
	}
}

// SetAddr sets the listening address for the server
func SetAddr(addr string) ApplyFunc {
	return func(s *http.Server) {
		s.Addr = addr
	}
}

// Option is the interface to apply server options
type Option interface {
	Apply(s *http.Server)
}

// ApplyFunc is a function type that satisfies the options interface.
type ApplyFunc func(s *http.Server)

func (a ApplyFunc) Apply(s *http.Server) {
	a(s)
}

// Options takes multiple options and applies them in one call
type Options []Option

func (o Options) Apply(s *http.Server) {
	for _, option := range o {
		option.Apply(s)
	}
}
