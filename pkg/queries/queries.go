package queries

import (
	"errors"
	"math/rand"
	"memnixrest/app/database"
	"memnixrest/app/models"
	"memnixrest/pkg/core"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DeleteRating(c *fiber.Ctx, user *models.User, deck *models.Deck) models.ResponseHTTP {
	db := database.DBConn

	rating := new(models.Rating)

	if err := db.Where("ratings.user_id = ? AND ratings.deck_id = ?", user.ID, deck.ID).Delete(&rating).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		}
	}
	return models.ResponseHTTP{
		Success: true,
		Message: "Rating deleted",
		Data:    *rating,
		Count:   1,
	}
}

func GenerateRating(c *fiber.Ctx, rating *models.Rating) models.ResponseHTTP {
	db := database.DBConn
	deck := new(models.Deck)

	if err := db.First(&deck, rating.DeckID).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Can't find deck matching the ID provided",
			Data:    nil,
			Count:   0,
		}
	}

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", rating.UserID, rating.DeckID).First(&access).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permissions to rate this deck",
			Data:    nil,
			Count:   0,
		}
	}

	if access.Permission < models.AccessStudent {
		return models.ResponseHTTP{
			Success: false,
			Message: "You are not subscribed to this deck",
			Data:    nil,
			Count:   0,
		}
	}

	oldRating := new(models.Rating)

	if err := db.Joins("User").Joins("Deck").Where("ratings.user_id = ? AND ratings.deck_id = ?", rating.UserID, rating.DeckID).First(&oldRating).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			db.Preload("User").Preload("Deck").Create(rating)
			oldRating = rating
		}
	} else {
		oldRating.Value = rating.Value
		db.Preload("User").Preload("Deck").Save(oldRating)
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Success rate the deck",
		Data:    *oldRating,
		Count:   1,
	}
}

// GenerateAdminAccess
func GenerateCreatorAccess(c *fiber.Ctx, user *models.User, deck *models.Deck) models.ResponseHTTP {
	db := database.DBConn

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", user.ID, deck.ID).Find(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			access.DeckID = deck.ID
			access.UserID = user.ID
			access.Permission = models.AccessOwner
			db.Create(access)
		}

	} else {
		if access.Permission >= models.AccessStudent {
			return models.ResponseHTTP{
				Success: false,
				Message: "You are already subscribed to this deck. You can't become an owner...",
				Data:    nil,
				Count:   0,
			}
		} else {
			access.DeckID = deck.ID
			access.UserID = user.ID
			access.Permission = models.AccessOwner
			db.Preload("User").Preload("Deck").Save(access)
		}
	}

	log := CreateLog(models.LogSubscribe, user.Username+" subscribed to "+deck.DeckName)
	_ = CreateUserLog(*user, *log)
	_ = CreateDeckLog(*deck, *log)

	return models.ResponseHTTP{
		Success: true,
		Message: "Success register a creator access !",
		Data:    *access,
		Count:   1,
	}
}

// GenerateAccess
func GenerateAccess(c *fiber.Ctx, user *models.User, deck *models.Deck) models.ResponseHTTP {
	db := database.DBConn

	if deck.Status != models.DeckPublic && user.Permissions != models.PermAdmin {
		return models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permissions to subscribe to this deck!",
			Data:    nil,
			Count:   0,
		}
	}

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", user.ID, deck.ID).Find(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			access.DeckID = deck.ID
			access.UserID = user.ID
			access.Permission = models.AccessStudent
			db.Preload("User").Preload("Deck").Create(access)
		}

	} else {
		if access.Permission >= models.AccessStudent {
			return models.ResponseHTTP{
				Success: false,
				Message: "You are already subscribed to this deck",
				Data:    nil,
				Count:   0,
			}
		} else {
			access.DeckID = deck.ID
			access.UserID = user.ID
			access.Permission = models.AccessStudent
			db.Preload("User").Preload("Deck").Save(access)
		}
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Success register an access",
		Data:    *access,
		Count:   1,
	}
}

func CheckAccess(c *fiber.Ctx, user *models.User, deck *models.Deck, perm models.AccessPermission) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", user.ID, deck.ID).First(&access).Error; err != nil {
		access.Permission = models.AccessNone
	}

	if access.Permission < perm {
		return models.ResponseHTTP{
			Success: false,
			Message: "You don't have the permission to access this deck!",
			Data:    *access,
			Count:   1,
		}
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Success checking access permissions",
		Data:    *access,
		Count:   1,
	}
}

func PostMem(c *fiber.Ctx, user models.User, card models.Card, validation models.CardResponseValidation) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	memDate := new(models.MemDate)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		&user.ID, card.ID).First(&memDate).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "MemDate Not Found",
			Data:    nil,
		}
	}

	ex_mem := FetchMem(c, memDate, &user)
	if ex_mem.Efactor == 0 {
		ex_mem = models.Mem{
			UserID:     user.ID,
			CardID:     card.ID,
			Quality:    0,
			Repetition: 0,
			Efactor:    2.5,
			Interval:   0,
		}
	}

	core.UpdateMem(c, &ex_mem, validation)

	return models.ResponseHTTP{
		Success: true,
		Message: "Success Post Mem",
		Data:    nil,
	}
}

func PopulateMemDate(c *fiber.Ctx, user *models.User, deck *models.Deck) models.ResponseHTTP {
	db := database.DBConn // DB Conn
	var cards []models.Card

	if err := db.Joins("Deck").Where("cards.deck_id = ?", deck.ID).Find(&cards).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: err.Error(),
			Data:    nil,
			Count:   0,
		}
	}

	ch := make(chan models.ResponseHTTP)

	for _, s := range cards {
		go func(c *fiber.Ctx, user *models.User, card models.Card) {
			res := GenerateMemDate(c, user, &card)
			ch <- res
		}(c, user, s)
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Success generated mem_date",
		Data:    nil,
		Count:   0,
	}

}

func GenerateMemDate(c *fiber.Ctx, user *models.User, card *models.Card) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	memDate := new(models.MemDate)

	if err := db.Joins("User").Joins("Card").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?", user.ID, card.ID).First(&memDate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			memDate = &models.MemDate{
				UserID:   user.ID,
				CardID:   card.ID,
				DeckID:   card.DeckID,
				NextDate: time.Now(),
			}

			db.Preload("User").Preload("Card").Create(memDate)
		} else {
			return models.ResponseHTTP{
				Success: false,
				Message: err.Error(),
				Data:    nil,
				Count:   0,
			}
		}
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Success generate MemDate",
		Data:    *memDate,
		Count:   1,
	}
}

// FetchAnswers
func FetchAnswers(c *fiber.Ctx, card *models.Card) []models.Answer {
	var answers []models.Answer
	db := database.DBConn // DB Conn

	if err := db.Joins("Card").Where("answers.card_id = ?", card.ID).Limit(3).Order("random()").Find(&answers).Error; err != nil {
		return nil
	}

	return answers
}

func FetchMem(c *fiber.Ctx, memDate *models.MemDate, user *models.User) models.Mem {
	db := database.DBConn // DB Conn

	mem := new(models.Mem)
	if err := db.Joins("Card").Where("mems.card_id = ? AND mems.user_id = ?", memDate.CardID, user.ID).Order("id desc").First(&mem).Error; err != nil {
		mem.Efactor = 0
	}
	return *mem
}

func GenerateAnswers(c *fiber.Ctx, memDate *models.MemDate) []string {
	var answersList []string

	res := FetchAnswers(c, &memDate.Card)

	if len(res) >= 3 {
		answersList = append(answersList, res[0].Answer, res[1].Answer, res[2].Answer, memDate.Card.Answer)

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(answersList), func(i, j int) { answersList[i], answersList[j] = answersList[j], answersList[i] })
	}

	return answersList

}

func FetchNextCard(c *fiber.Ctx, user *models.User) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	memDate := new(models.MemDate)
	var answersList []string

	// Get next card
	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		user.ID).Joins("Card").Joins("Deck").Where("mem_dates.user_id = ? AND accesses.permission >= ?",
		&user.ID, models.AccessStudent).Limit(1).Order("next_date asc").Find(&memDate).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next card not found",
			Data:    nil,
		}
	}

	mem := FetchMem(c, memDate, user)
	if mem.Efactor <= 1.4 || mem.Quality <= 1 || mem.Repetition < 2 {
		answersList = GenerateAnswers(c, memDate)
		if len(answersList) == 4 {
			memDate.Card.Type = 2 // MCQ
		}
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Get Next Card",
		Data: models.ResponseCard{
			Card:    memDate.Card,
			Answers: answersList,
		},
	}
}

func FetchNextCardByDeck(c *fiber.Ctx, user *models.User, deckID string) models.ResponseHTTP {
	db := database.DBConn // DB Conn

	memDate := new(models.MemDate)
	var answersList []string

	// Get next card
	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.deck_id = ?",
		&user.ID, deckID).Limit(1).Order("next_date asc").Find(&memDate).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next card by deck not found",
			Data:    nil,
		}
	}

	mem := FetchMem(c, memDate, user)
	if mem.Efactor <= 1.4 || mem.Quality <= 1 || mem.Repetition < 2 || memDate.Card.Type == 2 {
		answersList = GenerateAnswers(c, memDate)
		if len(answersList) == 4 {
			memDate.Card.Type = 2 // MCQ
		}
	}

	return models.ResponseHTTP{
		Success: true,
		Message: "Get Next Card By Deck",
		Data: models.ResponseCard{
			Card:    memDate.Card,
			Answers: answersList,
		},
	}
}

// FetchNextTodayCard
func FetchNextTodayCard(c *fiber.Ctx, user *models.User) models.ResponseHTTP {
	db := database.DBConn // DB Conn
	memDate := new(models.MemDate)

	var answersList []string

	// Get next card with date condition
	t := time.Now()

	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		user.ID).Joins("Card").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.next_date < ? AND accesses.permission >= ?",
		&user.ID, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour), models.AccessStudent).Limit(1).Order("next_date asc").Find(&memDate).Error; err != nil {
		return models.ResponseHTTP{
			Success: false,
			Message: "Next today card not found",
			Data:    nil,
		}
	}
	mem := FetchMem(c, memDate, user)
	if mem.Efactor <= 1.4 || mem.Quality <= 1 || mem.Repetition < 2 || memDate.Card.Type == 2 {
		answersList = GenerateAnswers(c, memDate)
		if len(answersList) == 4 {
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
