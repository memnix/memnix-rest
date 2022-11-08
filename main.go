package main

import (
	cryptoRand "crypto/rand"
	"encoding/binary"
	"fmt"
	_ "github.com/arsmn/fiber-swagger/v2"
	"github.com/memnix/memnixrest/app/auth"
	"github.com/memnix/memnixrest/app/routes"
	"github.com/memnix/memnixrest/data/infrastructures"
	"github.com/memnix/memnixrest/models"
	"github.com/memnix/memnixrest/pkg/queries"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math/rand"
)

func init() {
	var b [8]byte
	_, err := cryptoRand.Read(b[:])
	if err != nil {
		panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

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
	// Try to connect to the infrastructures
	if err := infrastructures.Connect(); err != nil {
		log.Panic("Can't connect infrastructures:", err.Error())
	}

	// Create cache session
	if err := infrastructures.CreateCache(); err != nil {
		log.Panic("Can't create cache session:", err.Error())
	}

	// Connect to RabbitMQ
	if _, err := infrastructures.Rabbit(); err != nil {
		log.Panic("Can't connect to rabbitMq: ", err)
	}

	// Init the secret key
	auth.Init()

	// Disconnect from RabbitMQ*
	defer func(conn *amqp.Connection) {
		_ = conn.Close()
		fmt.Println("Disconnected to RabbitMQ")
	}(infrastructures.RabbitMQ.Connection)

	// Close RabbitMQ channel
	defer func(ch *amqp.Channel) {
		_ = ch.Close()
	}(infrastructures.RabbitMQ.Channel)

	// Models to migrate
	var migrates []interface{}
	migrates = append(migrates, models.Access{}, models.Card{}, models.Deck{},
		models.User{}, models.Mem{}, models.Answer{}, models.MemDate{}, models.Mcq{})

	// AutoMigrate models
	for i := 0; i < len(migrates); i++ {
		err := infrastructures.GetDBConn().AutoMigrate(&migrates[i])
		if err != nil {
			log.Panic("Can't auto migrate models:", err.Error())
		}
	}

	// Create queries cache for the first time
	queries.InitCache()

	// Create the app
	app := routes.New()
	// Listen to port 1812
	log.Panic(app.Listen(":1813"))
}
