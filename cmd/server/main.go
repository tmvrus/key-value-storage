package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tmvrus/key-value-storage/internal/app"
	"github.com/tmvrus/key-value-storage/internal/config"
	"github.com/tmvrus/key-value-storage/internal/storage"
	"gopkg.in/yaml.v3"
)

const defaultConfigFile = "./config.yml"

func fillConfig(cfg *config.Config, fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	return nil
}

func initLogger(fileName, level string) *slog.Logger {
	w := os.Stdout
	f, err := os.Open(fileName)
	if err != nil {
		slog.Error("failed to open file for logging, use stdout", err.Error(), "file_name", fileName)
	} else {
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
	err := fillConfig(cfg, configFile)
	if err != nil {
		slog.Error("failed to fill config, use default values", "error", err.Error())
	}

	log := initLogger(cfg.Logging.Output, cfg.Logging.Level)

	err = app.
		New(cfg, storage.New(cfg), log).
		Run(cxt)

	if err != nil {
		log.Error("failed to run application", "error", err.Error())
		os.Exit(1)
	}

}
