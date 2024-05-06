package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/bytedance/gopkg/util/gctuner"
	v2 "github.com/memnix/memnix-rest/app/v2"
	"github.com/memnix/memnix-rest/cmd/v2/config"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/crypto"
	"github.com/memnix/memnix-rest/pkg/json"
	"github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/pkg/errors"
)

func main() {
	configPath := config.GetConfigPath(config.IsDevelopment())

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ Error loading config: %s", err.Error())
	}

	logger.GetLogger().SetLogLevel(cfg.Log.GetSlogLevel()).CreateGlobalHandler()

	setup(cfg)

	e := v2.CreateEchoInstance(cfg.Server)

	if cfg.Log.GetSlogLevel().Level() == slog.LevelDebug {
		slog.Debug("🔧 Debug mode enabled")
	}

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

func setup(cfg *config.Config) {
	setupCrypto(cfg)

	setupInfrastructures(cfg)

	gcTuning()

	setupJwt(cfg)

	setupOAuth(cfg)

	slog.Info("✅ setup completed!")
}

func shutdown() error {
	e := v2.GetEchoInstance()

	slog.Info("🔒 Server shutting down...")

	var shutdownError error

	if err := e.Shutdown(context.Background()); err != nil {
		// Add the error to the shutdownError
		shutdownError = errors.Wrap(shutdownError, err.Error())
	}

	slog.Info("🧹 Running cleanup tasks...")

	infrastructures.GetPgxConnInstance().ClosePgx()

	slog.Info("✅ Disconnected from database")

	err := infrastructures.GetRedisManagerInstance().CloseRedis()
	if err != nil {
		shutdownError = errors.Wrap(shutdownError, err.Error())
	} else {
		slog.Info("✅ Disconnected from Redis")
	}

	slog.Info("✅ Disconnected from Sentry")

	slog.Info("✅ Cleanup tasks completed!")

	return shutdownError
}

func setupJwt(cfg *config.Config) {
	// Parse the keys
	if err := crypto.GetKeyManagerInstance().ParseEd25519Key(); err != nil {
		log.Fatal("❌ Error parsing keys", slog.Any("error", err))
	}

	// Create the JWT instance
	jwtInstance := jwt.NewJWTInstance(cfg.Auth.JWTHeaderLen, cfg.Auth.JWTExpiration, crypto.GetKeyManagerInstance().GetPublicKey(), crypto.GetKeyManagerInstance().GetPrivateKey())

	jwt.GetJwtInstance().SetJwt(jwtInstance)

	slog.Info("✅ Created JWT instance")
}

func setupCrypto(cfg *config.Config) {
	crypto.GetCryptoHelperInstance().SetCryptoHelper(crypto.NewBcryptCrypto(cfg.Auth.Bcryptcost))

	slog.Info("✅ Created Crypto instance")
}

func setupOAuth(cfg *config.Config) {
	oauth.GetJSONHelperInstance().SetJSONHelper(json.NewJSON(&json.NativeJSON{}))

	oauthConfig := oauth.GlobalConfig{
		CallbackURL: cfg.Server.Host,
		FrontendURL: cfg.Server.FrontendURL,
	}

	oauth.SetOauthConfig(oauthConfig)

	oauth.InitGithub(cfg.Auth.Github)
	oauth.InitDiscord(cfg.Auth.Discord)

	slog.Info("✅ Created OAuth instance")
}

func setupInfrastructures(cfg *config.Config) {
	err := infrastructures.NewPgxConnInstance(infrastructures.PgxConfig{
		DSN: cfg.Pgx.DSN,
	}).ConnectPgx()
	if err != nil {
		log.Fatal("❌ Error connecting to database", slog.Any("error", err))
	}

	slog.Info("✅ Connected to database")

	// Redis connection
	err = infrastructures.NewRedisInstance(cfg.Redis).ConnectRedis()
	if err != nil {
		slog.Error("❌ Error connecting to Redis", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("✅ Connected to Redis")

	ristrettoInstance := infrastructures.CreateRistrettoInstance(cfg.Ristretto)

	if err = ristrettoInstance.CreateRistrettoCache(); err != nil {
		log.Fatal("❌ Error creating Ristretto cache", slog.Any("error", err))
	}
	slog.Info("✅ Created Ristretto cache")
}

func gcTuning() {
	var limit float64 = 4 * config.GCLimit
	// Set the GC threshold to 70% of the limit
	threshold := uint64(limit * config.GCThresholdPercent)

	gctuner.Tuning(threshold)

	slog.Info(fmt.Sprintf("🔧 GC Tuning - Limit: %.2f GB, Threshold: %d bytes, GC Percent: %d, Min GC Percent: %d, Max GC Percent: %d",
		limit/(config.GCLimit),
		threshold,
		gctuner.GetGCPercent(),
		gctuner.GetMinGCPercent(),
		gctuner.GetMaxGCPercent()))

	slog.Info("✅ GC Tuning completed!")
}
