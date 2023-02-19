package main

import (
	"github.com/joho/godotenv"
	"github.com/memnix/memnix-rest/app/http"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	logger.CreateLogger()

	err = infrastructures.ConnectEdgeDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	defer func() {
		err = infrastructures.CloseEdgeDB()
		if err != nil {
			log.Fatal().Err(err).Msg("Error closing database connection")
		}
	}()

	/*
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
	*/

	log.Debug().Msg("Starting server")
	// Create the app
	app := http.New()
	// Listen to port 1815
	log.Fatal().Err(app.Listen(":1815")).Msg("Error listening to port 1815")

}
