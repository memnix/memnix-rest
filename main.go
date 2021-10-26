package main

import (
	"log"
	"memnixrest/app/models"
	"memnixrest/app/database"
	"memnixrest/pkg/routes"
)

func main() {
	// Try to connect to the database
	if err := database.Connect(); err != nil {
		log.Panic("Can't connect database:", err.Error())
	}

	// AutoMigrate models
	database.DBConn.AutoMigrate(&models.Access{})
	database.DBConn.AutoMigrate(&models.Card{})
	database.DBConn.AutoMigrate(&models.Deck{})
	database.DBConn.AutoMigrate(&models.History{})
	database.DBConn.AutoMigrate(&models.Revision{})
	database.DBConn.AutoMigrate(&models.User{})
	database.DBConn.AutoMigrate(&models.Mem{})
	database.DBConn.AutoMigrate(&models.Answer{})

	// Create the app
	app := routes.New()
	// Listen to port 1812
	log.Fatal(app.Listen(":1813"))
}
