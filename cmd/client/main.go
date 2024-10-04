package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmvrus/key-value-storage/pkg/client"
)

const defaultConnectionAddr = "localhost:32230"

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var serverAddr string
	flag.StringVar(&serverAddr, "address", defaultConnectionAddr, "")
	flag.Parse()

	remoteAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		log.Error("failed to resolve server address", "address", serverAddr, "error", err.Error())
		os.Exit(1)
	}

	con, err := net.DialTCP("tcp", nil, remoteAddr)
	if err != nil {
		log.Error("failed to dial", "address", serverAddr, "error", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := con.Close(); err != nil {
			log.Error("failed to close connection", "error", err.Error())
		}
	}()

	client.
		NewClient(con, os.Stdin, os.Stdout, log).
		Start(ctx)
}
