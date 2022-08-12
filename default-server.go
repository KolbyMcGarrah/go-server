package goserver

import (
	"net/http"
	"time"
)

type Config struct {
	Addr                string
	HealthCheckEndpoint string
	ReadTimeout         time.Duration
	ReadHeaderTimeout   time.Duration
	WriteTimeout        time.Duration
	IdleTimeout         time.Duration
}

func NewDefaultConfig() *Config {
	return &Config{
		Addr:                ":8080",
		HealthCheckEndpoint: "/healthz",
		ReadTimeout:         30 * time.Second,
		ReadHeaderTimeout:   30 * time.Second,
		WriteTimeout:        30 * time.Second,
		IdleTimeout:         1 * time.Minute,
	}
}

func NewDefaultServer(config *Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              config.Addr,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
		Handler:           handler,
	}
}
