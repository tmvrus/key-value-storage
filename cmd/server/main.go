package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmvrus/key-value-storage/internal/config"
	"github.com/tmvrus/key-value-storage/internal/logger"
	"github.com/tmvrus/key-value-storage/internal/server"
	"github.com/tmvrus/key-value-storage/internal/storage"
)

const defaultConfigFile = "./config.yml"

func main() {
	cxt, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	var configFile string
	flag.StringVar(&configFile, "config", defaultConfigFile, "")
	flag.Parse()

	cfg := config.NewConfigWithDefaults()
	err := config.FillWithFile(cfg, configFile)
	if err != nil {
		slog.Error("failed to fill config, use default values", "error", err.Error())
	}

	log := logger.New(cfg.Logging.Output, cfg.Logging.Level)

	err = server.
		New(cfg, storage.New(cfg), log).
		Run(cxt)

	if err != nil {
		log.Error("failed to run application", "error", err.Error())
		os.Exit(1)
	}

}
