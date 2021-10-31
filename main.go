package main

import (
	"log"
	"memnixrest/app/database"
	"memnixrest/app/models"
	"memnixrest/pkg/routes"

	_ "github.com/arsmn/fiber-swagger/v2"
)

// @title Memnix
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:1813
// @BasePath /api
func main() {
	// Try to connect to the database
	if err := database.Connect(); err != nil {
		log.Panic("Can't connect database:", err.Error())
	}

	// AutoMigrate models
	database.DBConn.AutoMigrate(&models.Access{})
	database.DBConn.AutoMigrate(&models.Card{})
	database.DBConn.AutoMigrate(&models.Deck{})
	database.DBConn.AutoMigrate(&models.User{})
	database.DBConn.AutoMigrate(&models.Mem{})
	database.DBConn.AutoMigrate(&models.Answer{})
	database.DBConn.AutoMigrate(&models.MemDate{})
	database.DBConn.AutoMigrate(&models.DeckLogs{})
	database.DBConn.AutoMigrate(&models.UserLogs{})

	// Create the app
	app := routes.New()
	// Listen to port 1812
	log.Fatal(app.Listen(":1813"))
}
