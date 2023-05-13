package main

import (
	"github.com/bytedance/gopkg/util/gctuner"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/memnix/memnix-rest/app/http"
	"github.com/memnix/memnix-rest/app/meilisearch"
	"github.com/memnix/memnix-rest/app/misc"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/internal"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	// Create new relic app
	err = infrastructures.NewRelic()
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating new relic app")
	}

	// Create logger
	logger.CreateNewRelicLogger()

	err = infrastructures.ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	defer func() {
		err = infrastructures.DisconnectDB()
		if err != nil {
			log.Fatal().Err(err).Msg("Error disconnecting from database")
		}
	}()

	if !fiber.IsChild() {
		// Models to migrate
		migrates := []interface{}{
			// Add models here
			domain.User{}, domain.Deck{},
		}

		// AutoMigrate models
		for i := 0; i < len(migrates); i++ {
			err = infrastructures.GetDBConn().AutoMigrate(&migrates[i])
			if err != nil {
				log.Error().Err(err).Msg("Can't auto migrate models")
			}
		}
	}

	// Init oauth
	infrastructures.InitOauth()

	// Redis connection
	err = infrastructures.ConnectRedis()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to redis")
	}
	defer func() {
		err := infrastructures.CloseRedis()
		if err != nil {
			log.Fatal().Err(err).Msg("Error closing redis connection")
		}
	}()

	// Connect to influxDB
	err = infrastructures.ConnectInfluxDB(config.EnvHelper)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to influxDB")
	}

	defer func() {
		err = infrastructures.DisconnectInfluxDB()
		if err != nil {
			log.Fatal().Err(err).Msg("Error disconnecting from influxDB")
		}
	}()

	// Create logger workers
	go misc.CreateLogger()

	// Connect MeiliSearch
	infrastructures.ConnectMeiliSearch(config.EnvHelper)

	if !fiber.IsChild() {
		log.Debug().Msg("Starting server")

		// Init MeiliSearch
		err = meilisearch.InitMeiliSearch(internal.InitializeMeiliSearch())
		if err != nil {
			log.Error().Err(err).Msg("Can't init MeiliSearch")
		}
	}

	var limit float64 = 4 * 1024 * 1024 * 1024
	// Set the GC threshold to 70% of the limit
	threshold := uint64(limit * config.GCThresholdPercent)

	gctuner.Tuning(threshold)

	// Create the app
	app := http.New()
	// Listen to port 1815
	log.Fatal().Err(app.Listen(":1815")).Msg("Error listening to port 1815")

	log.Info().Msg("Server stopped")
}
