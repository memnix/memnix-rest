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

	var migrates []interface{}
	_ = append(migrates, models.Access{}, models.Card{}, models.Deck{},
		models.User{}, models.Mem{}, models.Answer{}, models.MemDate{}, models.DeckLogs{},
		models.CardLogs{}, models.UserLogs{}, models.Logs{})

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
