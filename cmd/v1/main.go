package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/gopkg/util/gctuner"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	http "github.com/memnix/memnix-rest/app/v1"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/crypto"
	"github.com/memnix/memnix-rest/pkg/json"
	myJwt "github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/memnix/memnix-rest/pkg/oauth"
)

var version = "development"

func main() {
	configPath := config.GetConfigPath()

	// Load the config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("❌ Error loading config: %s", err.Error())
	}

	logger.NewLogger().SetLogLevel(slog.LevelInfo)

	if config.IsProduction() {
		log.Info("🚀 Running in production mode 🚀")
		cfg.Server.AppVersion = version
		cfg.Sentry.Release = "memnix@" + version
	}

	// setup the environment variables
	setup(cfg)

	log.Info("starting server 🚀", slog.String("version", cfg.Server.AppVersion))

	// Create the app
	app := http.New()

	// Listen from a different goroutine
	go func() {
		if err = app.Listen(":1815"); err != nil {
			log.Error("error starting server", slog.Any("error", err))
			// exit with error
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received

	shutdown(app)

	log.Info("server stopped")
}

func shutdown(app *fiber.App) {
	log.Info("🔒 Server shutting down...")
	_ = app.Shutdown()

	log.Info("🧹 Running cleanup tasks...")

	err := infrastructures.GetDBConnInstance().DisconnectDB()
	if err != nil {
		log.Error("❌ Error closing database connection")
	} else {
		log.Info("✅ Disconnected from database")
	}

	err = infrastructures.GetRedisManagerInstance().CloseRedis()
	if err != nil {
		log.Error("❌ Error closing Redis connection")
	} else {
		log.Info("✅ Disconnected from Redis")
	}

	err = infrastructures.ShutdownTracer()
	if err != nil {
		log.Error("❌ Error closing Tracer connection")
	} else {
		log.Info("✅ Disconnected from Tracer")
	}

	sentry.Flush(config.SentryFlushTimeout)
	log.Info("✅ Disconnected from Sentry")

	log.Info("✅ Cleanup tasks completed!")
}

func setup(cfg *config.Config) {
	setupCrypto(cfg)

	setupInfrastructures(cfg)

	gcTuning()

	setupJwt(cfg)

	setupOAuth(cfg)

	if !fiber.IsChild() {
		// Migrate the models
		migrate()
	}

	log.Info("✅ setup completed!")
}

func setupJwt(cfg *config.Config) {
	// Parse the keys
	if err := config.GetKeyManagerInstance().ParseEd25519Key(); err != nil {
		log.Fatal("❌ Error parsing keys", slog.Any("error", err))
	}

	// Create the JWT instance
	jwtInstance := myJwt.NewJWTInstance(cfg.Auth.JWTHeaderLen, cfg.Auth.JWTExpiration, config.GetKeyManagerInstance().GetPublicKey(), config.GetKeyManagerInstance().GetPrivateKey())

	config.GetJwtInstance().SetJwt(jwtInstance)

	log.Info("✅ Created JWT instance")
}

func setupCrypto(cfg *config.Config) {
	crypto.GetCryptoHelperInstance().SetCryptoHelper(crypto.NewBcryptCrypto(cfg.Auth.Bcryptcost))

	log.Info("✅ Created Crypto instance")
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

	log.Info("✅ Created OAuth instance")
}

func setupInfrastructures(cfg *config.Config) {
	err := infrastructures.GetDBConnInstance().ConnectDB(cfg.Database.DSN)
	if err != nil {
		log.Fatal("❌ Error connecting to database", slog.Any("error", err))
	}

	log.Info("✅ Connected to database")

	// Redis connection
	err = infrastructures.GetRedisManagerInstance().ConnectRedis(cfg.Redis)
	if err != nil {
		log.Fatal("❌ Error connecting to Redis")
	}
	log.Info("✅ Connected to Redis")

	// Connect to the tracer
	err = infrastructures.InitTracer(cfg.Sentry)
	if err != nil {
		log.Fatal("❌ Error connecting to Tracer", slog.Any("error", err))
	}
	log.Info("✅ Created Tracer")

	if err = infrastructures.GetCacheInstance().CreateRistrettoCache(); err != nil {
		log.Fatal("❌ Error creating Ristretto cache", slog.Any("error", err))
	}
	log.Info("✅ Created Ristretto cache")
}

func gcTuning() {
	var limit float64 = 4 * config.GCLimit
	// Set the GC threshold to 70% of the limit
	threshold := uint64(limit * config.GCThresholdPercent)

	gctuner.Tuning(threshold)

	log.Info(fmt.Sprintf("🔧 GC Tuning - Limit: %.2f GB, Threshold: %d bytes, GC Percent: %d, Min GC Percent: %d, Max GC Percent: %d",
		limit/(config.GCLimit),
		threshold,
		gctuner.GetGCPercent(),
		gctuner.GetMinGCPercent(),
		gctuner.GetMaxGCPercent()))

	log.Info("✅ GC Tuning completed!")
}

func migrate() {
	// Models to migrate
	migrates := []domain.Model{
		&domain.User{}, &domain.Card{}, &domain.Deck{}, &domain.Mcq{},
	}

	log.Info("⚙️ Starting database migration...")

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		step := i + 1
		err := infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			log.Error(fmt.Sprintf("❌ Error migrating model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		} else {
			log.Info(fmt.Sprintf("✅ Migration completed for model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		}
	}

	log.Info("✅ Database migration completed!")
}
