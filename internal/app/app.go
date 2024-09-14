package app

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/tmvrus/key-value-storage/internal/compute/parser"
	"github.com/tmvrus/key-value-storage/internal/domain"
)

type App struct {
	input  io.Reader
	output io.Writer

	s storage
}

func New(reader io.Reader, writer io.Writer, s storage) App {
	return App{input: reader, output: writer, s: s}
}

func (a App) write(s string) {
	_, err := a.output.Write([]byte(s + "\n"))
	if err != nil {
		// TODO: log here
	}
}

func (a App) error(err error) {
	msg := "ERROR: " + err.Error()
	a.write(msg)
}

func (a App) doCmd(ctx context.Context, c domain.Command) (string, error) {
	switch c.Type {
	case domain.CommandGet:
		return a.s.Get(ctx, c.Key)
	case domain.CommandDelete:
		return "", a.s.Delete(ctx, c.Key)
	case domain.CommandSet:
		return "", a.s.Set(ctx, c.Key, c.Value)
	default:
		return "", fmt.Errorf("invalid cmd type: %q", c.Type)
	}
}

func (a App) Run(ctx context.Context) {
	input := bufio.NewScanner(a.input)
	for input.Scan() {

		s := input.Text()

		cmd, err := parser.Parse(s)
		if err != nil {
			a.error(err)
			continue
		}

		res, err := a.doCmd(ctx, cmd)
		if err != nil {
			a.error(err)
			continue
		}

		if res == "" {
			res = "OK"
		}
		a.write(res)
	}
}
