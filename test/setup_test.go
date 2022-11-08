package test

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/memnix/memnixrest/app/routes"
	infrastructures2 "github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/pkg/models"
	"github.com/memnix/memnixrest/pkg/queries"
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
			if _, got := Setup(); got != nil {
				t.Errorf("Setup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Setup() (*fiber.App, error) {
	// Try to connect to the infrastructures
	if err := infrastructures2.Connect(); err != nil {
		log.Panic("Can't connect infrastructures:", err.Error())
	}

	// Connect to RabbitMQ
	if _, err := infrastructures2.Rabbit(); err != nil {
		log.Panic("Can't connect to rabbitMq: ", err)
	}

	// Disconnect from RabbitMQ*
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
		fmt.Println("Disconnected to RabbitMQ")
	}(infrastructures2.RabbitMqConn)

	// Close RabbitMQ channel
	defer func(ch *amqp.Channel) {
		_ = ch.Close()
	}(infrastructures2.RabbitMqChan)

	// Models to migrate
	var migrates []interface{}
	migrates = append(migrates, models.Access{}, models.Card{}, models.Deck{},
		models.User{}, models.Mem{}, models.Answer{}, models.MemDate{}, models.Mcq{})

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		err := infrastructures2.DBConn.AutoMigrate(&migrates[i])
		if err != nil {
			log.Panic("Can't auto migrate models:", err.Error())
		}
	}

	// Create the app
	app := routes.New()

	queries.InitCache()

	return app, nil
}
