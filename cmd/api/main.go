package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/gopkg/util/gctuner"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnix-rest/app/http"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/crypto"
	myJwt "github.com/memnix/memnix-rest/pkg/jwt"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/memnix/memnix-rest/pkg/oauth"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func main() {
	// setup the logger
	zapLogger, undo := logger.CreateZapLogger()

	configPath := config.GetConfigPath()

	// Load the config
	cfg, err := config.UseConfig(configPath)
	if err != nil {
		zapLogger.Fatal("❌ Error loading config", zap.Error(err))
	}
	// setup the environment variables
	setup(cfg)

	zapLogger.Info("starting server")

	// Create the app
	app := http.New()

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":1815"); err != nil {
			zapLogger.Panic("error starting server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received

	shutdown(app)

	zapLogger.Info("server stopped")

	if err := zapLogger.Sync(); err != nil {
		return // can't even log, just exit
	}
	undo()
}

func shutdown(app *fiber.App) {
	otelzap.L().Info("🔒 Server shutting down...")
	_ = app.Shutdown()

	otelzap.L().Info("🧹 Running cleanup tasks...")

	err := infrastructures.DisconnectDB()
	if err != nil {
		otelzap.L().Error("❌ Error closing database connection")
	} else {
		otelzap.L().Info("✅ Disconnected from database")
	}

	err = infrastructures.CloseRedis()
	if err != nil {
		otelzap.L().Error("❌ Error closing Redis connection")
	} else {
		otelzap.L().Info("✅ Disconnected from Redis")
	}

	err = infrastructures.ShutdownTracer()
	if err != nil {
		otelzap.L().Error("❌ Error closing Tracer connection")
	} else {
		otelzap.L().Info("✅ Disconnected from Tracer")
	}
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

	otelzap.L().Info("✅ setup completed!")
}

func setupJwt(cfg *config.Config) {
	// Parse the private key
	if err := config.ParseEd25519PrivateKey(); err != nil {
		otelzap.L().Fatal("❌ Error parsing private key", zap.Error(err))
	}

	// Parse the public key
	if err := config.ParseEd25519PublicKey(); err != nil {
		otelzap.L().Fatal("❌ Error parsing public key", zap.Error(err))
	}

	// Create the JWT instance
	jwtInstance := myJwt.NewJWTInstance(cfg.Auth.JWTHeaderLen, cfg.Auth.JWTExpiration, config.GetEd25519PublicKey(), config.GetEd25519PrivateKey())

	config.JwtInstance = jwtInstance

	otelzap.L().Info("✅ Created JWT instance")
}

func setupCrypto(cfg *config.Config) {
	crypto.InitCrypto(crypto.NewBcryptCrypto(cfg.Auth.Bcryptcost))

	otelzap.L().Info("✅ Created Crypto instance")
}

func setupOAuth(cfg *config.Config) {
	oauth.SetJSONHelper(config.JSONHelper)

	oauthConfig := oauth.GlobalConfig{
		CallbackURL: cfg.Server.Host,
		FrontendURL: cfg.Server.FrontendURL,
	}

	oauth.SetOauthConfig(oauthConfig)

	oauth.InitGithub(cfg.Auth.Github)
	oauth.InitDiscord(cfg.Auth.Discord)

	otelzap.L().Info("✅ Created OAuth instance")
}

func setupInfrastructures(cfg *config.Config) {
	err := infrastructures.ConnectDB(cfg.Database.DSN)
	if err != nil {
		otelzap.L().Fatal("❌ Error connecting to database", zap.Error(err))
	} else {
		otelzap.L().Info("✅ Connected to database")
	}

	// Redis connection
	err = infrastructures.ConnectRedis(cfg.Redis)
	if err != nil {
		otelzap.L().Fatal("❌ Error connecting to Redis")
	} else {
		otelzap.L().Info("✅ Connected to Redis")
	}

	// Connect to the tracer
	err = infrastructures.InitTracer(cfg.Tracing)
	if err != nil {
		otelzap.L().Fatal("❌ Error connecting to Tracer", zap.Error(err))
	} else {
		otelzap.L().Info("✅ Created Tracer")
	}

	if err = infrastructures.CreateRistrettoCache(); err != nil {
		otelzap.L().Fatal("❌ Error creating Ristretto cache", zap.Error(err))
	} else {
		otelzap.L().Info("✅ Created Ristretto cache")
	}

}

func gcTuning() {
	var limit float64 = 4 * config.GCLimit
	// Set the GC threshold to 70% of the limit
	threshold := uint64(limit * config.GCThresholdPercent)

	gctuner.Tuning(threshold)

	otelzap.L().Info(fmt.Sprintf("🔧 GC Tuning - Limit: %.2f GB, Threshold: %d bytes, GC Percent: %d, Min GC Percent: %d, Max GC Percent: %d",
		limit/(config.GCLimit),
		threshold,
		gctuner.GetGCPercent(),
		gctuner.GetMinGCPercent(),
		gctuner.GetMaxGCPercent()))

	otelzap.L().Info("✅ GC Tuning completed!")
}

func migrate() {
	// Models to migrate
	migrates := []domain.Model{
		&domain.User{}, &domain.Card{}, &domain.Deck{}, &domain.Mcq{},
	}

	otelzap.L().Info("⚙️ Starting database migration...")

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		step := i + 1
		err := infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			otelzap.L().Error(fmt.Sprintf("❌ Error migrating model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		} else {
			otelzap.L().Info(fmt.Sprintf("✅ Migration completed for model %s %d/%d", migrates[i].TableName(), step, len(migrates)))
		}
	}

	otelzap.L().Info("✅ Database migration completed!")
}
