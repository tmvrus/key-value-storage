package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmvrus/key-value-storage/internal/app"
	"github.com/tmvrus/key-value-storage/internal/storage/engine/inmemory"
)

func main() {
	cxt, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	app.
		New(os.Stdin, os.Stdout, inmemory.New()).
		Run(cxt)
}
