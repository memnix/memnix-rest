package queries

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
)

func CreateUserLog(user *models.User, logType models.UserLogType, message string) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	userLog := &models.UserLogs{
		UserID:  user.ID,
		LogType: logType,
		Message: message,
	}

	db.Preload("User").Preload("Card").Create(userLog)
	return models.ResponseHTTP{
		Success: true,
		Message: "Created a new user log entry",
		Data:    userLog,
		Count:   1,
	}
}

func CreateDeckLog(deck *models.Deck, logType models.DeckLogType, message string) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	deckLog := &models.DeckLogs{
		DeckID:  deck.ID,
		LogType: logType,
		Message: message,
	}

	db.Preload("User").Preload("Card").Create(deckLog)
	return models.ResponseHTTP{
		Success: true,
		Message: "Created a deck user log entry",
		Data:    deckLog,
		Count:   1,
	}
}

// TODO: Get last log by type
// TODO: Fully implement logs
