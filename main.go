package main

import (
	"fmt"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/routes"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"

	_ "github.com/arsmn/fiber-swagger/v2"
)

// @title Memnix
// @version 1.0
// @description Memnix API documentation
// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name memnix-jwt
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host https://api-memnix.yumenetwork.net
// @BasePath /api
func main() {
	// Try to connect to the database
	if err := database.Connect(); err != nil {
		log.Panic("Can't connect database:", err.Error())
	}

	if _, err := database.Rabbit(); err != nil {
		log.Panic("Can't connect to rabbitMq: ", err)
	}

	defer func(conn *amqp.Connection) {
		_ = conn.Close()
		fmt.Println("Disconnected to RabbitMQ")

	}(database.RabbitMqConn)

	defer func(ch *amqp.Channel) {
		_ = ch.Close()
	}(database.RabbitMqChan)

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
	// Listen to port 1812
	log.Fatal(app.Listen(":1813"))
}
