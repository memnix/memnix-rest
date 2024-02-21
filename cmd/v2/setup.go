package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/bytedance/gopkg/util/gctuner"
	"github.com/getsentry/sentry-go"
	v2 "github.com/memnix/memnix-rest/app/v2"
	"github.com/memnix/memnix-rest/cmd/v2/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/crypto"
	"github.com/memnix/memnix-rest/pkg/json"
	"github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/pkg/errors"
)

func setup(cfg *config.Config) {
	setupCrypto(cfg)

	setupInfrastructures(cfg)

	gcTuning()

	setupJwt(cfg)

	setupOAuth(cfg)

	migrate()

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

	err := infrastructures.GetDBConnInstance().DisconnectDB()
	if err != nil {
		shutdownError = errors.Wrap(shutdownError, err.Error())
	} else {
		slog.Info("✅ Disconnected from database")
	}

	err = infrastructures.GetRedisManagerInstance().CloseRedis()
	if err != nil {
		shutdownError = errors.Wrap(shutdownError, err.Error())
	} else {
		slog.Info("✅ Disconnected from Redis")
	}

	sentry.Flush(config.SentryFlushTimeout)

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
	err := infrastructures.NewDBConnInstance(infrastructures.DatabaseConfig{
		DSN:             cfg.Database.DSN,
		SQLMaxIdleConns: cfg.Database.SQLMaxIdleConns,
		SQLMaxOpenConns: cfg.Database.SQLMaxOpenConns,
	}).ConnectDB()
	if err != nil {
		log.Fatal("❌ Error connecting to database", slog.Any("error", err))
	}

	slog.Info("✅ Connected to database")

	// Redis connection
	err = infrastructures.NewRedisInstance(cfg.Redis).ConnectRedis()
	if err != nil {
		log.Fatal("❌ Error connecting to Redis")
	}
	slog.Info("✅ Connected to Redis")

	// Connect to the tracer
	err = infrastructures.NewTracerInstance(infrastructures.SentryConfig{
		DSN:                cfg.Sentry.DSN,
		Environment:        cfg.Sentry.Environment,
		Debug:              cfg.Sentry.Debug,
		Release:            cfg.Sentry.Release,
		TracesSampleRate:   cfg.Sentry.TracesSampleRate,
		ProfilesSampleRate: cfg.Sentry.ProfilesSampleRate,
		Name:               "fiber-rest",
		WithStacktrace:     false,
	}).ConnectTracer()
	if err != nil {
		log.Fatal("❌ Error connecting to Tracer", slog.Any("error", err))
	}
	slog.Info("✅ Created Tracer")

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

func migrate() {
	// Models to migrate
	migrates := []domain.Model{
		&domain.User{}, &domain.Card{}, &domain.Deck{}, &domain.Mcq{},
	}

	slog.Info("⚙️ Starting database migration...")

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		step := i + 1
		err := infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			slog.Error(fmt.Sprintf("❌ Error migrating model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		} else {
			slog.Info(fmt.Sprintf("✅ Migration completed for model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		}
	}

	slog.Info("✅ Database migration completed!")
}
