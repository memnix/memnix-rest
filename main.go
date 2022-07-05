package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/routes"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"

	_ "github.com/arsmn/fiber-swagger/v2"
)

// @title Memnix
// @version 1.0
// @description Memnix API
// @securityDefinitions.apikey Beaver
// @in header
// @name Authorization
// @securityDefinitions.apikey Admin
// @in header
// @name Authorization
// @termsOfService https://github.com/memnix/memnix/blob/main/PRIVACY.md
// @contact.name API Support
// @contact.email contact@memnix.app
// @license.name BSD 3-Clause License
// @license.url https://github.com/memnix/memnix-rest/blob/main/LICENSE
// @host http://192.168.1.151:1813/
// @BasePath /v1
func main() {
	app := Setup() // Create the app

	// Listen to port 1813
	log.Fatal(app.Listen(":1813"))

}

func Setup() *fiber.App {
	// Try to connect to the database
	if err := database.Connect(); err != nil {
		log.Panic("Can't connect database:", err.Error())
	}

	// Connect to RabbitMQ
	if _, err := database.Rabbit(); err != nil {
		log.Panic("Can't connect to rabbitMq: ", err)
	}

	// Disconnect from RabbitMQ*
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
		fmt.Println("Disconnected to RabbitMQ")
	}(database.RabbitMqConn)

	// Close RabbitMQ channel
	defer func(ch *amqp.Channel) {
		_ = ch.Close()
	}(database.RabbitMqChan)

	// Models to migrate
	var migrates []interface{}
	migrates = append(migrates, models.Access{}, models.Card{}, models.Deck{},
		models.User{}, models.Mem{}, models.Answer{}, models.MemDate{}, models.Mcq{})

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		err := database.DBConn.AutoMigrate(&migrates[i])
		if err != nil {
			log.Panic("Can't auto migrate models:", err.Error())
		}
	}

	// Create the app
	app := routes.New()

	return app
}
