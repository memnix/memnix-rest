package queries

import (
	"math/rand"
	"memnixrest/app/database"
	"memnixrest/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

// FetchNextTodayMemByUserAndDeck
func FetchNextTodayCard(c *fiber.Ctx, user *models.User) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	mem := new(models.Mem)
	memDate := new(models.MemDate)
	var answers []models.Answer

	// Get next card with date condition
	t := time.Now()

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.next_date < ?",
		&user.ID, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour)).Limit(1).Order("next_date asc").Find(&memDate).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next today card not found",
			Data:    nil,
		}
	}

	if err := db.Joins("Card").Where("mems.card_id = ? AND mems.user_id = ?", memDate.CardID, user.ID).Find(&mem).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next today card not found",
			Data:    nil,
		}
	}

	if mem.Total < 2 || mem.Efactor <= 1.4 || mem.Quality <= 1 || mem.Repetition < 2 {
		// Answers
	}

	if err := db.Joins("Card").Where("answers.card_id = ?", memDate.CardID).Limit(3).Order("random()").Find(&answers).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next today card not found",
			Data:    nil,
		}
	}

	var answersList []string
	if len(answers) >= 3 {
		answersList = append(answersList, answers[0].Answer, answers[1].Answer, answers[2].Answer, memDate.Card.Answer)

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(answersList), func(i, j int) { answersList[i], answersList[j] = answersList[j], answersList[i] })
	}

	responseCard := models.ResponseCard{
		Card:    memDate.Card,
		Answers: answersList,
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Get Next Today Card",
		Data:    responseCard,
	}
}
