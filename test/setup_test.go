package test

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/routes"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"testing"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name string
		want *fiber.App
	}{
		{
			name: "TestSetup",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Setup(); got != nil {
				t.Errorf("Setup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Setup() (error, *fiber.App) {
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

	return nil, app
}
