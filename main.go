package main

import (
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/memnix/memnix-rest/app/http"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	logger.CreateLogger()

	err = infrastructures.ConnectSentry()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing sentry")
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)

	// Connect to database
	log.Debug().Msg("Connecting to database")

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

	// Connect to redis
	err = infrastructures.ConnectRedis()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to redis")
	}
	defer func() {
		err = infrastructures.CloseRedis()
		if err != nil {
			log.Fatal().Err(err).Msg("Error closing redis connection")
		}
	}()

	// Models to migrate
	migrates := []interface{}{
		// Add models here
		domain.User{},
	}

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		err = infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			log.Error().Err(err).Msg("Can't auto migrate models")
		}
	}

	log.Debug().Msg("Starting server")
	// Create the app
	app := http.New()
	// Listen to port 1815
	log.Fatal().Err(app.Listen(":1815")).Msg("Error listening to port 1815")
}
