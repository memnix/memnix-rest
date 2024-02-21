package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	v2 "github.com/memnix/memnix-rest/app/v2"
	"github.com/memnix/memnix-rest/cmd/v2/config"
	"github.com/memnix/memnix-rest/pkg/logger"
)

func main() {
	configPath := config.GetConfigPath()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ Error loading config: %s", err.Error())
	}

	logger.NewLogger().SetLogLevel(slog.LevelInfo)

	e := v2.GetEchoInstance()

	setup(cfg)

	slog.Info("starting server 🚀", slog.String("version", cfg.Server.AppVersion))

	go func() {
		if err = e.Start(":3000"); err != nil {
			slog.Error("error starting server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	slog.Info("shutting down server")

	if err = shutdown(); err != nil {
		slog.Error("error shutting down server", slog.Any("error", err))
	}

	slog.Info("server stopped")
}
