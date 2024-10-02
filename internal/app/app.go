package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/tmvrus/key-value-storage/internal/config"
)

type App struct {
	log     *slog.Logger
	storage storage
	cfg     *config.Config

	sessionLimiter chan struct{}
}

func New(cfg *config.Config, s storage, l *slog.Logger) App {
	return App{
		log:            l,
		storage:        s,
		cfg:            cfg,
		sessionLimiter: make(chan struct{}, cfg.Network.MaxConnections),
	}
}

func (a App) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", a.cfg.Network.Address)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	a.log.Debug("ready to accept connections", "address", a.cfg.Network.Address)

	go func() {
		<-ctx.Done()
		if err := l.Close(); err != nil {
			a.log.Error("failed to close listener", "error", err.Error())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			a.log.Debug("got context done, stop application")
			return ctx.Err()
		default:

		}

		conn, err := l.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				a.log.Error("failed to accept connection", "error", err.Error())
			}
			continue
		}

		select {
		case a.sessionLimiter <- struct{}{}:
			a.log.Debug("start session", "src", conn.RemoteAddr().String())

		default:
			a.log.Debug("drop session due the limit", "src", conn.RemoteAddr().String())
			if err := conn.Close(); err != nil {
				a.log.Error("failed to close connection", "error", err.Error())
			}
			continue
		}

		go func() {
			defer func() {
				if err := conn.Close(); err != nil {
					a.log.Error("failed to close connection", "error", err.Error())
				}
			}()

			start := time.Now()
			cfg := sessionConfig{
				timeout:    a.cfg.Network.IdleTimeout,
				bufferSize: a.cfg.Network.MaxMessageSize.Int(),
			}
			newSession(a.log, a.storage, conn, cfg).start(ctx)
			<-a.sessionLimiter

			a.log.Debug("session finished", "src", conn.RemoteAddr().String(), "duration", time.Since(start).String())
		}()
	}
}
