package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

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

	e := v2.CreateEchoInstance(cfg.Server)

	setup(cfg)

	slog.Info("starting server 🚀", slog.String("version", cfg.Server.AppVersion))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		if err = e.Start(); err != nil {
			slog.Error("error starting server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	const shutdownTimeout = 10 * time.Second

	<-ctx.Done()
	_, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	slog.Info("shutting down server")

	if err = shutdown(); err != nil {
		slog.Error("error shutting down server", slog.Any("error", err))
	}

	slog.Info("server stopped")
}
