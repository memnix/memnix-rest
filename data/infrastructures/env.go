package infrastructures

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func loadDBVar() dbParams {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConfig := new(dbParams)

	if os.Getenv("APP_ENV") == "dev" {
		log.Println("Running in development mode")
		dbConfig.User = os.Getenv("DEBUG_DB_USER")         // Get DB_USER from env
		dbConfig.Password = os.Getenv("DEBUG_DB_PASSWORD") // Get DB_PASSWORD from env
		dbConfig.Host = os.Getenv("DEBUG_DB_HOST")         // Get DB_HOST from env
		dbConfig.DBName = os.Getenv("DEBUG_DB_DB")         // Get DB_DB (db name) from env
		dbConfig.Port = os.Getenv("DEBUG_DB_PORT")         // Get DB_PORT from env
	} else {
		log.Println("Running in production mode")
		dbConfig.User = os.Getenv("DB_USER")         // Get DB_USER from env
		dbConfig.Password = os.Getenv("DB_PASSWORD") // Get DB_PASSWORD from env
		dbConfig.Host = os.Getenv("DB_HOST")         // Get DB_HOST from env
		dbConfig.DBName = os.Getenv("DB_DB")         // Get DB_DB (db name) from env
		dbConfig.Port = os.Getenv("DB_PORT")         // Get DB_PORT from env
	}

	return *dbConfig
}

func loadRabbitMQVar() RabbitMq {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rabbitMQConfig := new(RabbitMq)

	if os.Getenv("APP_ENV") == "dev" {
		log.Println("Running in development mode")
		rabbitMQConfig.rabbitMqUrl = os.Getenv("DEBUG_RABBIT_MQ") // Get DB_HOST from env
	} else {
		log.Println("Running in production mode")
		rabbitMQConfig.rabbitMqUrl = os.Getenv("RABBIT_MQ") // Get DB_HOST from env
	}

	return *rabbitMQConfig
}
