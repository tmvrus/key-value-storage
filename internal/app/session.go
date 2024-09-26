package app

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"

	"github.com/tmvrus/key-value-storage/internal/compute/parser"
	"github.com/tmvrus/key-value-storage/internal/domain"
)

type session struct {
	log     *slog.Logger
	storage storage
	socket  socket
}

func newSession(l *slog.Logger, st storage, s socket) session {
	return session{
		log:     l,
		storage: st,
		socket:  s,
	}
}

func (a session) start(ctx context.Context) {
	input := bufio.NewScanner(a.socket)
	for input.Scan() {

		text := input.Text()
		if err := input.Err(); err != nil {
			a.log.Error("got error while read socket", "error", err.Error())
			continue
		}

		cmd, err := parser.Parse(text)
		if err != nil {
			a.writeError(err)
			continue
		}

		res, err := a.doCmd(ctx, cmd)
		if err != nil {
			a.writeError(err)
			continue
		}

		if res == "" {
			res = "OK"
		}
		a.writeStringLn(res)
	}
}

func (a session) writeStringLn(s string) {
	_, err := a.socket.Write([]byte(s + "\n"))
	if err != nil {
		a.log.Error("failed to write socket", "error", err.Error())
	}
}

func (a session) writeError(err error) {
	msg := "ERROR: " + err.Error()
	a.writeStringLn(msg)
}

func (a session) doCmd(ctx context.Context, c domain.Command) (string, error) {
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
