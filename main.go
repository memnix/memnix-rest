package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bytedance/gopkg/util/gctuner"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/memnix/memnix-rest/app/http"
	"github.com/memnix/memnix-rest/app/meilisearch"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup the environment variables
	setupEnv()

	// Setup the garbage collector
	gcTuning()

	// Setup the infrastructures
	setupInfrastructures()

	if !fiber.IsChild() {
		// Migrate the models
		migrate()

		// Init MeiliSearch
		err := meilisearch.InitMeiliSearch(internal.InitializeMeiliSearch())
		if err != nil {
			log.Error().Err(err).Msg("Can't init MeiliSearch")
		}
	}

	log.Debug().Msg("Starting server...")

	// Create the app
	app := http.New()

	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":1815"); err != nil {
			log.Panic().Err(err).Msg("Error listening to port 1815")
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received

	shutdown(app)

	log.Info().Msg("Server stopped")
}

func shutdown(app *fiber.App) {
	log.Info().Msg("🔒 Server shutting down...")
	_ = app.Shutdown()

	log.Info().Msg("🧹 Running cleanup tasks...")

	err := infrastructures.DisconnectDB()
	if err != nil {
		log.Error().Err(err).Msg("❌ Error disconnecting from database")
	} else {
		log.Info().Msg("✅ Disconnected from database")
	}

	err = infrastructures.CloseRedis()
	if err != nil {
		log.Error().Err(err).Msg("❌ Error closing Redis connection")
	} else {
		log.Info().Msg("✅ Disconnected from Redis")
	}

	err = infrastructures.ShutdownTracer()
	if err != nil {
		log.Error().Err(err).Msg("❌ Error closing Tracer")
	} else {
		log.Info().Msg("✅ Disconnected from Tracer")
	}
}

func migrate() {
	// Models to migrate
	migrates := []domain.Model{
		&domain.User{}, &domain.Card{}, &domain.Deck{}, &domain.Mcq{},
	}

	log.Info().Msg("⚙️ Starting database migration...")

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		step := i + 1
		err := infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			log.Error().Err(err).Int("model", step).Msg("❌ Can't auto migrate models")
		} else {
			log.Info().Msgf("✅ Migration completed for model %s %d/%d", migrates[i].TableName(), step, len(migrates))
		}
	}

	log.Info().Msg("✅ Database migration completed!")
}

func setupEnv() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	logger.CreateLogger()

	// Init oauth
	infrastructures.InitOauth()
}

func setupInfrastructures() {
	err := infrastructures.ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Error connecting to database")
	} else {
		log.Info().Msg("✅ Connected to database")
	}

	// Redis connection
	err = infrastructures.ConnectRedis()
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Error connecting to Redis")
	} else {
		log.Info().Msg("✅ Connected to Redis")
	}

	// Connect MeiliSearch
	err = infrastructures.ConnectMeiliSearch(config.EnvHelper)
	if err != nil {
		log.Error().Err(err).Msg("❌ Error connecting to MeiliSearch")
	} else {
		log.Info().Msg("✅ Connected to MeiliSearch")
	}

	// Connect to the tracer
	err = infrastructures.InitTracer()
	if err != nil {
		log.Error().Err(err).Msg("❌ Error connecting to the tracer")
	} else {
		log.Info().Msg("✅ Connected to the tracer")
	}
}

func gcTuning() {
	var limit float64 = 4 * config.GCLimit
	// Set the GC threshold to 70% of the limit
	threshold := uint64(limit * config.GCThresholdPercent)

	gctuner.Tuning(threshold)

	log.Debug().
		Msgf("🔧 GC Tuning - Limit: %.2f GB, Threshold: %d bytes, GC Percent: %d, Min GC Percent: %d, Max GC Percent: %d",
			limit/(config.GCLimit),
			threshold,
			gctuner.GetGCPercent(),
			gctuner.GetMinGCPercent(),
			gctuner.GetMaxGCPercent())

	log.Debug().Msg("✅ GC Tuning completed.")
}
