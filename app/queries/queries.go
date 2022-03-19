package queries

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/memnix/memnixrest/app/models"
	"github.com/memnix/memnixrest/pkg/core"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/utils"
	"gorm.io/gorm"
)

// FillResponseDeck returns a filled models.ResponseDeck
// This function might become a method of models.ResponseDeck
func FillResponseDeck(deck *models.Deck, permission models.AccessPermission) models.ResponseDeck {
	db := database.DBConn

	deckResponse := new(models.ResponseDeck)

	deckResponse.Deck = *deck
	deckResponse.DeckID = deck.ID
	deckResponse.Permission = permission

	if owner := deck.GetOwner(); owner.ID != 0 {
		publicUser := new(models.PublicUser)

		publicUser.Set(&owner)

		deckResponse.Owner = *publicUser
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

// GenerateCreatorAccess sets an user as a deck creator
func GenerateCreatorAccess(user *models.User, deck *models.Deck) *models.ResponseHTTP {
	db := database.DBConn
	// TODO: Change models.User & models.Deck to uint
	access := new(models.Access)
	res := new(models.ResponseHTTP)

	if err := db.Where("accesses.user_id = ? AND accesses.deck_id =?", user.ID, deck.ID).Find(&access).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			access.Set(user.ID, deck.ID, models.AccessOwner)
			db.Create(access)
		} else {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(deck.ID)
		fmt.Println(user.ID)
		res.GenerateError(utils.ErrorForbidden)
		return res
	}

	res.GenerateSuccess("Success register a creator access !", *access, 1)
	return res
}

// GenerateAccess sets a default student access to a deck for a given user
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
			access.Set(user.ID, deck.ID, models.AccessStudent)
			db.Preload("User").Preload("Deck").Create(access)
		}

	} else {
		if access.Permission >= models.AccessStudent {
			res.GenerateError(utils.ErrorAlreadySub)
			return res

		} else {
			access.Set(user.ID, deck.ID, models.AccessStudent)
			db.Preload("User").Preload("Deck").Save(access)
		}
	}

	res.GenerateSuccess("Success register an access", *access, 1)
	return res
}

// CheckAccess verifies if a given user as the right models.Permission to perform an action on a deck
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

// CheckCardLimit verifies that a deck can handle more cards
func CheckCardLimit(permission models.Permission, deckID uint) bool {
	db := database.DBConn // DB Conn
	var count int64

	if err := db.Table("cards").Where("cards.deck_id = ?", deckID).Count(&count).Error; err != nil {
		//TODO: Handle error
		return true
	}

	if permission < models.PermMod && count >= utils.MaximumCardDeck {
		return false
	}

	return true
}

// CheckDeckLimit verifies that the user hasn't reached the limit
func CheckDeckLimit(user *models.User) bool {
	db := database.DBConn // DB Conn
	var count int64

	if err := db.Table("accesses").Where("accesses.user_id = ? AND accesses.permission = ?", user.ID, models.AccessOwner).Count(&count).Error; err != nil {
		//TODO: Handle error
		return true
	}

	if user.Permissions < models.PermMod && count >= utils.MaximumDeckNormalUser {
		return false
	}

	return true
}

// PostMem updates MemDate & Mem
func PostMem(user *models.User, card *models.Card, validation *models.CardResponseValidation, training bool) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	//TODO: Replace struct params with ids
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

// PopulateMemDate with default value for a given user & deck
// This is used on deck sub
func PopulateMemDate(user *models.User, deck *models.Deck) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	var cards []models.Card
	res := new(models.ResponseHTTP)

	if err := db.Joins("Deck").Where("cards.deck_id = ?", deck.ID).Find(&cards).Error; err != nil {
		res.GenerateError(err.Error()) // MemDate not found
		return res
	}

	for i := range cards {
		_ = GenerateMemDate(user.ID, cards[i].ID, cards[i].DeckID)
	}
	res.GenerateSuccess("Success generated mem_date", nil, 0)
	return res
}

// GetSubUsers returns a list of users sub to a deck
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

// GenerateMemDate with default nextDate
func GenerateMemDate(userID, cardID, deckID uint) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	res := new(models.ResponseHTTP)

	memDate := new(models.MemDate)

	if err := db.Joins("User").Joins("Card").Where("mem_dates.user_id = ? AND mem_dates.card_id = ?", userID, cardID).First(&memDate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			memDate.SetDefaultNextDate(userID, cardID, deckID)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			mem.Efactor = 0
		}
	}
	return *mem
}

func GenerateMCQ(memDate *models.MemDate, userID uint) []string {

	mem := FetchMem(memDate.CardID, userID)

	var answersList []string
	if mem.IsMCQ() || memDate.Card.Type == models.CardMCQ {

		answersList = memDate.Card.GetMCQAnswers()
		if len(answersList) == 4 {
			memDate.Card.Type = models.CardMCQ // MCQ
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

	if err := db.Joins("Deck").Joins("Card").Where("mem_dates.deck_id = ? AND mem_dates.user_id = ?", deckID, userID).Find(&memDates).Error; err != nil {
		res.GenerateError(err.Error())
		return res
	}
	responseCard := new(models.ResponseCard)
	var answersList []string

	for i := range memDates {

		answersList = GenerateMCQ(&memDates[i], userID)
		responseCard.Set(&memDates[i].Card, answersList)

		result = append(result, *responseCard)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })

	res.GenerateSuccess("Success getting next card", result, len(result))
	return res

}

func FetchTodayCard(userID uint) *models.ResponseHTTP {
	db := database.DBConn // DB Conn
	t := time.Now()

	res := new(models.ResponseHTTP)
	var memDates []models.MemDate

	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		userID).Joins("Card").Joins("Deck").Where("mem_dates.user_id = ? AND mem_dates.next_date < ? AND accesses.permission >= ?",
		userID, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour), models.AccessStudent).Order("next_date asc").Find(&memDates).Error; err != nil {
		res.GenerateError("Today's memDate not found")
		return res
	}

	var answersList []string
	var responseCardList []models.ResponseCard
	responseCard := new(models.ResponseCard)

	for index := range memDates {
		answersList = GenerateMCQ(&memDates[index], userID)
		responseCard.Set(&memDates[index].Card, answersList)
		responseCardList = append(responseCardList, *responseCard)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(responseCardList), func(i, j int) { responseCardList[i], responseCardList[j] = responseCardList[j], responseCardList[i] })

	res.GenerateSuccess("Success getting next today's cards", responseCardList, len(responseCardList))
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

	responseCard.Set(&memDate.Card, answersList)

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
	responseCard.Set(&memDate.Card, answersList)

	res.GenerateSuccess("Success getting next card", responseCard, 1)
	return res
}
