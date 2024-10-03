package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tmvrus/key-value-storage/internal/config"
	"github.com/tmvrus/key-value-storage/internal/server"
	"github.com/tmvrus/key-value-storage/internal/storage"
)

const defaultConfigFile = "./config.yml"

func initLogger(fileName, level string) *slog.Logger {
	w := os.Stdout
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend|os.ModeExclusive)
	if err != nil {
		slog.Error("failed to open file for logging, use stdout", err.Error(), "file_name", fileName)
	} else {
		go func() {
			t := time.NewTicker(time.Second)
			for range t.C {
				if err := f.Sync(); err != nil {
					slog.Error("failed to sync log file", "file_name", fileName, "error", err.Error())
				}
			}
		}()

		w = f
	}

	m := map[string]slog.Level{
		"info":  slog.LevelInfo,
		"debug": slog.LevelDebug,
		"error": slog.LevelError,
	}
	l, ok := m[level]
	if !ok {
		slog.Error("failed to find log level, use debug", "level", level)
		l = slog.LevelDebug
	}

	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: l}))
}

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

	log := initLogger(cfg.Logging.Output, cfg.Logging.Level)

	err = server.
		New(cfg, storage.New(cfg), log).
		Run(cxt)

	if err != nil {
		log.Error("failed to run application", "error", err.Error())
		os.Exit(1)
	}

}
