package queries

import (
	"errors"
	"math/rand"
	"memnixrest/app/database"
	"memnixrest/app/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// FetchAnswers
func FetchAnswers(c *fiber.Ctx, card *models.Card) []models.Answer {
	var answers []models.Answer
	db := database.DBConn // DB Conn

	if err := db.Joins("Card").Where("answers.card_id = ?", card.ID).Limit(3).Order("random()").Find(&answers).Error; err != nil {
		return nil
	}

	return answers
}

// FetchNextTodayCard
func FetchNextTodayCard(c *fiber.Ctx, user *models.User) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	mem := new(models.Mem)
	memDate := new(models.MemDate)
	var answersList []string

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

	if err := db.Joins("Card").Where("mems.card_id = ? AND mems.user_id = ?", memDate.CardID, user.ID).Order("id asc").First(&mem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			mem = nil
		} else {
			return models.ResponseHTTP{
				Success: false,
				Message: "Next today card not found",
				Data:    nil,
			}
		}
	}

	if mem == nil || mem.Efactor <= 1.4 || mem.Quality <= 1 || mem.Repetition < 2 {
		// Answers
		res := FetchAnswers(c, &memDate.Card)

		if len(res) >= 3 {
			answersList = append(answersList, res[0].Answer, res[1].Answer, res[2].Answer, memDate.Card.Answer)

			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(answersList), func(i, j int) { answersList[i], answersList[j] = answersList[j], answersList[i] })

			memDate.Card.Type = 2 // MCQ
		}
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Get Next Today Card",
		Data: models.ResponseCard{
			Card:    memDate.Card,
			Answers: answersList,
		},
	}
}
