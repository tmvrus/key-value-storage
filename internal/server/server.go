package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/tmvrus/key-value-storage/internal/config"
)

type Server struct {
	log     *slog.Logger
	storage storage
	cfg     *config.Config

	sessionLimiter chan struct{}
}

func New(cfg *config.Config, s storage, l *slog.Logger) Server {
	return Server{
		log:            l,
		storage:        s,
		cfg:            cfg,
		sessionLimiter: make(chan struct{}, cfg.Network.MaxConnections),
	}
}

func (s Server) Run(ctx context.Context) error {
	l, err := net.Listen("tcp", s.cfg.Network.Address)
	if err != nil {
		return fmt.Errorf("net listen: %w", err)
	}

	s.log.Debug("ready to accept connections", "address", s.cfg.Network.Address)

	go func() {
		<-ctx.Done()
		if err := l.Close(); err != nil {
			s.log.Error("failed to close listener", "error", err.Error())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			s.log.Debug("got context done, stop application")
			return ctx.Err()
		default:

		}

		conn, err := l.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				s.log.Error("failed to accept connection", "error", err.Error())
			}
			continue
		}

		select {
		case s.sessionLimiter <- struct{}{}:
			s.log.Debug("start session", "src", conn.RemoteAddr().String())

		default:
			s.log.Debug("drop session due the limit", "src", conn.RemoteAddr().String())
			if err := conn.Close(); err != nil {
				s.log.Error("failed to close connection", "error", err.Error())
			}
			continue
		}

		go func() {
			defer func() {
				if err := conn.Close(); err != nil {
					s.log.Error("failed to close connection", "error", err.Error())
				}
			}()

			start := time.Now()
			cfg := handlerConfig{
				timeout:    s.cfg.Network.IdleTimeout,
				bufferSize: s.cfg.Network.MaxMessageSize.Int(),
			}
			newHandler(s.log, s.storage, conn, cfg).startHandling(ctx)
			<-s.sessionLimiter

			s.log.Debug("session finished", "src", conn.RemoteAddr().String(), "duration", time.Since(start).String())
		}()
	}
}
