package queries

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
)

func CreateLog(logType models.LogType, message string) models.Logs {
	db := database.DBConn // DB Conn

	log := models.Logs{
		LogType: logType,
		Message: message,
	}

	db.Create(log)
	return log

}

func CreateUserLog(user models.User, log models.Logs) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	userLog := models.UserLogs{
		UserID: user.ID,
		LogID:  log.ID,
	}

	db.Preload("User").Preload("Card").Create(userLog)
	return models.ResponseHTTP{
		Success: true,
		Message: "Created a new user log entry",
		Data:    userLog,
		Count:   1,
	}
}

func CreateDeckLog(deck models.Deck, log models.Logs) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	deckLog := models.DeckLogs{
		DeckID: deck.ID,
		LogID:  log.ID,
	}

	db.Preload("User").Preload("Card").Create(deckLog)
	return models.ResponseHTTP{
		Success: true,
		Message: "Created a deck user log entry",
		Data:    deckLog,
		Count:   1,
	}
}
