package queries

import (
	"memnixrest/app/database"
	"memnixrest/app/models"
)

func CreateLog(logType models.LogType, message string) *models.Logs {
	db := database.DBConn // DB Conn

	log := &models.Logs{
		LogType: logType,
		Message: message,
	}

	db.Create(log)
	return log

}

func CreateUserLog(userID uint, log models.Logs) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	userLog := &models.UserLogs{
		UserID: userID,
		LogID:  log.ID,
	}

	db.Create(userLog)
	return models.ResponseHTTP{
		Success: true,
		Message: "Created an user log entry",
		Data:    *userLog,
		Count:   1,
	}
}

func CreateDeckLog(deckID uint, log models.Logs) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	deckLog := &models.DeckLogs{
		DeckID: deckID,
		LogID:  log.ID,
	}

	db.Create(deckLog)
	return models.ResponseHTTP{
		Success: true,
		Message: "Created a deck log entry",
		Data:    *deckLog,
		Count:   1,
	}
}

func CreateCardLog(cardID uint, log models.Logs) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	cardLog := &models.CardLogs{
		CardID: cardID,
		LogID:  log.ID,
	}

	db.Create(cardLog)
	return models.ResponseHTTP{
		Success: true,
		Message: "Created a card log entry",
		Data:    *cardLog,
		Count:   1,
	}
}
