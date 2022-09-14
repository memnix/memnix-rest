package queries

import (
	"errors"
	"fmt"
	"github.com/memnix/memnixrest/pkg/logger"
	"math/rand"
	"sort"
	"time"

	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/models"
	"github.com/memnix/memnixrest/pkg/utils"
	"gorm.io/gorm"
)

// UpdateSubUsers generates MemDate for sub users
func UpdateSubUsers(card *models.Card, user *models.User) error {
	var users []models.User
	var result *models.ResponseHTTP

	if result = GetSubUsers(card.DeckID); !result.Success {
		log := logger.CreateLog(fmt.Sprintf("Error from %s on deck %d - CreateNewCard: %s", user.Email, card.DeckID, result.Message),
			logger.LogQueryGetError).SetType(logger.LogTypeError).AttachIDs(user.ID, card.DeckID, card.ID)
		_ = log.SendLog()
		return errors.New("couldn't get sub users")
	}

	switch result.Data.(type) {
	default:
		return errors.New("couldn't get sub users")
	case []models.User:
		users = result.Data.([]models.User)
	}

	for i := range users {
		_ = GenerateMemDate(users[i].ID, card.ID, card.DeckID)
	}

	return nil
}

// FillResponseDeck returns a filled models.ResponseDeck
// This function might become a method of models.ResponseDeck
func FillResponseDeck(deck *models.Deck, permission models.AccessPermission, toggleToday bool) models.ResponseDeck {
	db := database.DBConn

	deckResponse := models.ResponseDeck{
		Deck:        *deck,
		DeckID:      deck.ID,
		Permission:  permission,
		ToggleToday: toggleToday,
		OwnerID:     0,
		Owner:       models.PublicUser{},
	}

	if owner := deck.GetOwner(); owner.ID != 0 {
		publicUser := new(models.PublicUser)

		publicUser.Set(&owner)

		deckResponse.Owner = *publicUser
		deckResponse.OwnerID = owner.ID
	}

	var count int64
	if err := db.Table("cards").Where("cards.deck_id = ?", deck.ID).Count(&count).Error; err != nil {
		deckResponse.CardCount = 0
	} else {
		deckResponse.CardCount = uint16(count)
	}
	return deckResponse
}

// GenerateCreatorAccess sets an user as a deck creator
func GenerateCreatorAccess(user *models.User, deck *models.Deck) *models.ResponseHTTP {
	db := database.DBConn
	// TODO: Change models.User & models.Deck to uint
	access := new(models.Access)
	res := new(models.ResponseHTTP)

	access.Set(user.ID, deck.ID, models.AccessOwner)
	db.Create(access)

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
		}
		access.Set(user.ID, deck.ID, models.AccessStudent)
		db.Preload("User").Preload("Deck").Save(access)
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

	if err := db.Table("cards").Where("cards.deck_id = ? AND cards.deleted_at IS NULL", deckID).Count(&count).Error; err != nil {
		//TODO: Handle error
		return true
	}

	if permission < models.PermMod && count >= utils.MaxCardDeck {
		return false
	}

	return true
}

// CheckCode prevents deck code from being duplicated
func CheckCode(key, code string) bool {
	db := database.DBConn // DB Conn
	var count int64

	if err := db.Table("decks").Where("decks.key = ? AND decks.code = ? AND decks.deleted_at IS NULL", key, code).Count(&count).Error; err != nil {
		// TODO: Handle error
		return true
	}

	if count != 0 {
		return false
	}

	return true
}

// CheckDeckLimit verifies that the user hasn't reached the limit
func CheckDeckLimit(user *models.User) bool {
	db := database.DBConn // DB Conn
	var count int64

	if err := db.Table("accesses").Where("accesses.user_id = ? AND accesses.permission = ? AND accesses.deleted_at IS NULL", user.ID, models.AccessOwner).Count(&count).Error; err != nil {
		//TODO: Handle error
		return true
	}

	if user.Permissions < models.PermMod && count >= utils.MaxDeckNormalUser {
		return false
	}

	return true
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

// GenerateMCQ returns a list of answer
func GenerateMCQ(memDate *models.MemDate, userID uint) []string {
	mem := FetchMem(memDate.CardID, userID)

	answersList := make([]string, 4)
	if mem.IsMCQ() || memDate.Card.Type == models.CardMCQ {
		answersList = memDate.Card.GetMCQAnswers()
		if len(answersList) == 4 {
			memDate.Card.Type = models.CardMCQ // MCQ
		}

		return answersList
	}

	return answersList
}

// FetchTrainingCards returns training cards
func FetchTrainingCards(userID, deckID uint) *models.ResponseHTTP {
	res := new(models.ResponseHTTP)
	db := database.DBConn // DB Conn

	var memDates []models.MemDate

	if err := db.Joins("Deck").Joins("Card").Where("mem_dates.deck_id = ? AND mem_dates.user_id = ?", deckID, userID).Find(&memDates).Error; err != nil {
		res.GenerateError(err.Error())
		return res
	}
	responseCard := new(models.ResponseCard)
	var answersList []string

	result := make([]models.ResponseCard, len(memDates))

	for i := range memDates {
		answersList = GenerateMCQ(&memDates[i], userID)
		responseCard.Set(&memDates[i], answersList)
		result[i] = *responseCard
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(result), func(i, j int) { result[i], result[j] = result[j], result[i] })

	res.GenerateSuccess("Success getting next card", result, len(result))
	return res
}

// FetchTodayCard return today cards
func FetchTodayCard(userID uint) *models.ResponseHTTP {
	db := database.DBConn // DB Conn

	res := new(models.ResponseHTTP)

	memDates, err := FetchTodayMemDate(userID)
	if err != nil {
		res.GenerateError(err.Error())
		return res
	}

	m, err := GenerateResponseCardMap(memDates, userID)
	if err != nil {
		res.GenerateError(err.Error())
		return res
	}

	todayResponse := new(models.TodayResponse)

	for key := range m {
		deck := new(models.Deck)
		_ = db.First(&deck, key).Error
		deckResponse := models.DeckResponse{
			DeckID: key,
			Cards:  m[key],
			Count:  len(m[key]),
			Deck:   *deck,
		}
		todayResponse.AppendDeckResponse(deckResponse)
	}

	sort.Slice(todayResponse.DecksReponses, func(i, j int) bool {
		return todayResponse.DecksReponses[i].Count < todayResponse.DecksReponses[j].Count
	})

	todayResponse.Count = len(todayResponse.DecksReponses)

	res.GenerateSuccess("Success getting next today's cards", todayResponse, len(memDates))
	return res
}
