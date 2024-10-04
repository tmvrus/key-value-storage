package client

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"syscall"
)

type Client struct {
	sender     sender
	interactor interactor
	log        *slog.Logger
}

func (c *Client) Start(ctx context.Context) {
	for {
		cmd, err := c.interactor.ReadCommand()
		if err != nil {
			c.log.Error("read new command", "error", err.Error())
			if criticalError(err) {
				return
			}

			continue
		}

		res, err := c.sender.Send(ctx, cmd)
		if err != nil {
			c.log.Error("send command", "error", err.Error())
			if criticalError(err) {
				return
			}

			continue
		}

		err = c.interactor.WriteResult(res)
		if err != nil {
			c.log.Error("send command", "error", err.Error())
			if criticalError(err) {
				return
			}
		}
	}
}

func NewClient(conn readerWriter, input io.Reader, output io.Writer, l *slog.Logger) *Client {
	di := struct {
		io.Reader
		io.Writer
	}{input, output}

	return &Client{
		sender:     newSender(conn),
		interactor: newConsoleInteractor(di),
		log:        l,
	}
}

func criticalError(err error) bool {
	return errors.Is(err, io.EOF) ||
		errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded)
}
