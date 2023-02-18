package main

import (
	"github.com/joho/godotenv"
	"github.com/memnix/memnix-rest/app/http"
	"github.com/memnix/memnix-rest/domain"
	"github.com/memnix/memnix-rest/infrastructures"
	"github.com/memnix/memnix-rest/pkg/logger"
	"github.com/rs/zerolog/log"
)

func main() {
	logger.CreateLogger()

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	err = infrastructures.ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
	}
	defer func() {
		err = infrastructures.CloseDB()
		if err != nil {
			log.Fatal().Err(err).Msg("Error closing database connection")
		}
	}()

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
	var migrates []interface{}
	migrates = append(migrates, domain.Barista{})

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		err = infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			log.Panic().Err(err).Msg("Can't auto migrate models")
		}
	}

	// Create the app
	app := http.New()
	// Listen to port 1812
	log.Fatal().Err(app.Listen(":1815")).Msg("Error listening to port 1815")

}
