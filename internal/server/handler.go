package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
	"time"

	"github.com/tmvrus/key-value-storage/internal/compute/parser"
	"github.com/tmvrus/key-value-storage/internal/domain"
)

type handlerConfig struct {
	timeout    time.Duration
	bufferSize int
}

type handler struct {
	log     *slog.Logger
	storage storage
	conn    socket
	cfg     handlerConfig
}

func newHandler(l *slog.Logger, st storage, s socket, cfg handlerConfig) handler {
	return handler{
		log:     l,
		storage: st,
		conn:    s,
		cfg:     cfg,
	}
}

func (a handler) startHandling(ctx context.Context) {
	input := bufio.NewScanner(a.conn)
	input.Buffer(make([]byte, a.cfg.bufferSize), a.cfg.bufferSize)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := a.conn.SetReadDeadline(time.Now().Add(a.cfg.timeout)); err != nil {
			a.handleError(err, "set read deadline")
			return
		}

		if !input.Scan() {
			a.handleError(input.Err(), "scan")
			return
		}

		text := input.Text()
		if err := input.Err(); err != nil {
			if a.handleError(err, "read command") || errors.Is(err, os.ErrDeadlineExceeded) {
				return
			}
			continue
		}

		cmd, err := parser.Parse(text)
		if err != nil {
			if a.handleError(a.writeError(err), "parse command") {
				return
			}
			continue
		}

		res, err := a.doCmd(ctx, cmd)
		if err != nil {
			if a.handleError(a.writeError(err), "exec cmd") {
				return
			}
			continue
		}

		err = a.writeResult(res)
		if a.handleError(err, "write result") {
			return
		}
	}
}

func (a handler) writeResult(res string) error {
	if res == "" {
		res = "OK"
	}

	err := a.writeStringLn(res)
	if err != nil {
		return fmt.Errorf("write string: %w", err)
	}

	return nil
}

func (a handler) handleError(err error, message string) (stop bool) {
	if err == nil {
		return
	}

	a.log.Error("got handler error", "error", fmt.Errorf("%s: %w", message, err))

	stop = errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
	return
}

func (a handler) writeStringLn(s string) error {
	err := a.conn.SetWriteDeadline(time.Now().Add(a.cfg.timeout))
	if err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	_, err = a.conn.Write([]byte(s + "\n"))
	if err != nil {
		return fmt.Errorf("failed to write conn: %w", err)
	}
	return nil
}

func (a handler) writeError(err error) error {
	msg := "ERROR: " + err.Error()
	return a.writeStringLn(msg)
}

func (a handler) doCmd(ctx context.Context, c domain.Command) (string, error) {
	switch c.Type {
	case domain.CommandGet:
		return a.storage.Get(ctx, c.Key)
	case domain.CommandDelete:
		return "", a.storage.Delete(ctx, c.Key)
	case domain.CommandSet:
		return "", a.storage.Set(ctx, c.Key, c.Value)
	default:
		return "", fmt.Errorf("invalid cmd type: %q", c.Type)
	}
}
