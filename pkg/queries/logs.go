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

func CreateUserLog(userID uint, log *models.Logs) *models.ResponseHTTP {
	db := database.DBConn // DB Conn

	res := new(models.ResponseHTTP)

	userLog := &models.UserLogs{
		UserID: userID,
		LogID:  log.ID,
	}

	db.Create(userLog)

	res.GenerateSuccess("Created a user log entry", *userLog, 1)
	return res
}

func CreateDeckLog(deckID uint, log *models.Logs) *models.ResponseHTTP {
	db := database.DBConn // DB Conn

	res := new(models.ResponseHTTP)

	deckLog := &models.DeckLogs{
		DeckID: deckID,
		LogID:  log.ID,
	}

	db.Create(deckLog)

	res.GenerateSuccess("Created a deck log entry", *deckLog, 1)
	return res
}

func CreateCardLog(cardID uint, log *models.Logs) *models.ResponseHTTP {
	db := database.DBConn // DB Conn

	res := new(models.ResponseHTTP)

	cardLog := &models.CardLogs{
		CardID: cardID,
		LogID:  log.ID,
	}

	db.Create(cardLog)
	res.GenerateSuccess("Created a card log entry", *cardLog, 1)
	return res
}
