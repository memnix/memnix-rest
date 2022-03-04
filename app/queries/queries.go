package queries

import (
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"memnixrest/app/models"
	"memnixrest/pkg/core"
	"memnixrest/pkg/database"
	"memnixrest/pkg/utils"
	"time"
)

func FillResponseDeck(deck *models.Deck, permission models.AccessPermission) models.ResponseDeck {
	db := database.DBConn

	deckResponse := new(models.ResponseDeck)

	deckResponse.Deck = *deck
	deckResponse.DeckID = deck.ID
	deckResponse.Permission = permission

	if owner := deck.GetOwner(); owner.ID != 0 {
		deckResponse.Owner = owner
		deckResponse.OwnerId = owner.ID
	}

	var count int64
	if err := db.Table("cards").Where("cards.deck_id = ?", deck.ID).Count(&count).Error; err != nil {
		deckResponse.CardCount = 0
	} else {
		deckResponse.CardCount = count
	}
	return *deckResponse
}

// GenerateCreatorAccess function
func GenerateCreatorAccess(user *models.User, deck *models.Deck) *models.ResponseHTTP {
	db := database.DBConn

	access := new(models.Access)
	res := new(models.ResponseHTTP)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", user.ID, deck.ID).Find(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			access.Fill(user.ID, deck.ID, models.AccessOwner)
			db.Create(access)
		}
	} else {
		res.GenerateError(utils.ErrorForbidden)
		return res
	}

	log := CreateLog(models.LogSubscribe, user.Username+" subscribed to "+deck.DeckName)
	_ = CreateUserLog(user.ID, log)
	_ = CreateDeckLog(deck.ID, log)

	res.GenerateSuccess("Success register a creator access !", *access, 1)
	return res
}

// GenerateAccess function
func GenerateAccess(user *models.User, deck *models.Deck) *models.ResponseHTTP {
	db := database.DBConn
	res := new(models.ResponseHTTP)

	if deck.Status != models.DeckPublic && user.Permissions != models.PermAdmin {
		res.GenerateError(utils.ErrorForbidden)
		return res
	}

	access := new(models.Access)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id =?", user.ID, deck.ID).Find(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			access.Fill(user.ID, deck.ID, models.AccessStudent)
			db.Preload("User").Preload("Deck").Create(access)
		}

	} else {
		if access.Permission >= models.AccessStudent {
			res.GenerateError(utils.ErrorAlreadySub)
			return res

		} else {
			access.Fill(user.ID, deck.ID, models.AccessStudent)
			db.Preload("User").Preload("Deck").Save(access)
		}
	}

	res.GenerateSuccess("Success register an access", *access, 1)
	return res
}

func CheckAccess(userID, deckID uint, perm models.AccessPermission) *models.ResponseHTTP {
	db := database.DBConn // DB Conn

	access := new(models.Access)
	res := new(models.ResponseHTTP)

	if err := db.Joins("User").Joins("Deck").Where("accesses.user_id = ? AND accesses.deck_id = ?", userID, deckID).First(&access).Error; err != nil {
		access.Permission = models.AccessNone
	}

	if access.Permission < perm {
		res.GenerateError(utils.ErrorForbidden)
		return res
	}

	res.GenerateSuccess("Success checking access permissions", *access, 1)
	return res
}

func PostMem(user *models.User, card *models.Card, validation *models.CardResponseValidation, training bool) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(models.ResponseHTTP)

	memDate := new(models.MemDate)

	if err := db.Joins("Card").Joins("User").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?",
		user.ID, card.ID).First(&memDate).Error; err != nil {
		res.GenerateError(utils.ErrorRequestFailed) // MemDate not found
		// TODO: Create a default MemDate
		return res
	}

	exMem := FetchMem(memDate.CardID, user.ID)
	if exMem.Efactor == 0 {
		exMem.FillDefaultValues(user.ID, card.ID)
	}

	if training {
		core.UpdateMemTraining(&exMem, validation.Validate)
	} else {
		core.UpdateMem(&exMem, validation.Validate)
	}
	res.GenerateSuccess("Success Post Mem", nil, 0)
	return res
}

func PopulateMemDate(user *models.User, deck *models.Deck) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	var cards []models.Card
	res := new(models.ResponseHTTP)

	if err := db.Joins("Deck").Where("cards.deck_id = ?", deck.ID).Find(&cards).Error; err != nil {
		res.GenerateError(err.Error()) // MemDate not found
		return res
	}

	for _, s := range cards {
		_ = GenerateMemDate(user.ID, s.ID, s.DeckID)
	}
	res.GenerateSuccess("Success generated mem_date", nil, 0)
	return res
}

func GetSubUsers(deckID uint) *models.ResponseHTTP {
	res := new(models.ResponseHTTP)

	db := database.DBConn // DB Conn
	var users []models.User

	if err := db.Joins("left join accesses ON users.id = accesses.user_id AND accesses.deck_id = ?", deckID).Where("accesses.permission > ?", models.AccessNone).Find(&users).Error; err != nil {
		res.GenerateError(err.Error())
		return res
	}
	res.GenerateSuccess("Success getting sub users", users, len(users))
	return res
}

func GenerateMemDate(userID, cardID, deckID uint) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(models.ResponseHTTP)

	memDate := new(models.MemDate)

	if err := db.Joins("User").Joins("Card").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?", userID, cardID).First(&memDate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			memDate.Generate(userID, cardID, deckID)
			db.Create(memDate)
		} else {
			res.GenerateError(err.Error())
			return res
		}
	}
	res.GenerateSuccess("Success generate MemDate", memDate, 1)
	return res
}

func FetchMem(cardID, userID uint) models.Mem {
	db := database.DBConn // DB Conn

	mem := new(models.Mem)
	if err := db.Joins("Card").Where("mems.card_id = ? AND mems.user_id = ?", cardID, userID).Order("id desc").First(&mem).Error; err != nil {
		mem.Efactor = 0
	}
	return *mem
}

func GenerateMCQ(memDate *models.MemDate, userID uint) []string {

	mem := FetchMem(memDate.Card.ID, userID)

	var answersList []string
	if mem.IsMCQ() || memDate.Card.Type == models.CardMCQ {

		answersList = memDate.Card.GetMCQAnswers()
		if len(answersList) == 4 {
			memDate.Card.Type = 2 // MCQ
		}
		return answersList
	}

	return answersList
}

func FetchTrainingCards(userID, deckID uint) *models.ResponseHTTP {
	res := new(models.ResponseHTTP)
	db := database.DBConn // DB Conn
	var result []models.ResponseCard

	var memDates []models.MemDate

	if err := db.Joins("Deck").Where("mem_dates.deck_id = ? AND mem_dates.user_id = ?", deckID, userID).Find(&memDates).Error; err != nil {
		res.GenerateError(err.Error())
		return res
	}
	responseCard := new(models.ResponseCard)
	var answersList []string

	for i := range memDates {

		answersList = GenerateMCQ(&memDates[i], userID)
		responseCard.Generate(memDates[i].Card, answersList)

		result = append(result, *responseCard)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })

	res.GenerateSuccess("Success getting next card", result, len(result))
	return res

}

func FetchNextTodayCard(userID uint) *models.ResponseHTTP {
	res := new(models.ResponseHTTP)
	responseCard := new(models.ResponseCard)
	memDate := new(models.MemDate)
	var answersList []string

	if result := memDate.GetNextToday(userID); !result.Success {
		res.GenerateError("Next card not found")
		return res
	}
	answersList = GenerateMCQ(memDate, userID)

	responseCard.Generate(memDate.Card, answersList)

	res.GenerateSuccess("Success getting next card", responseCard, 1)
	return res
}

func FetchNextCard(userID, deckID uint) *models.ResponseHTTP {
	res := new(models.ResponseHTTP)
	responseCard := new(models.ResponseCard)
	memDate := new(models.MemDate)
	var answersList []string

	if deckID != 0 {
		if result := memDate.GetNextByDeck(userID, deckID); !result.Success {
			res.GenerateError("Next card not found")
			return res
		}
	} else {
		if result := memDate.GetNext(userID); !result.Success {
			res.GenerateError("Next card not found")
			return res
		}

	}

	answersList = GenerateMCQ(memDate, userID)
	responseCard.Generate(memDate.Card, answersList)

	res.GenerateSuccess("Success getting next card", responseCard, 1)
	return res
}
