package queries

import (
	"errors"
	"github.com/memnix/memnixrest/cache"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/models"
	"gorm.io/gorm"
	"sync"
	"time"
)

var memDateCache *cache.Cache

func InitCache() {
	memDateCache = cache.NewCache()
}

func GetCache() *cache.Cache {
	return memDateCache
}

func ClearCache() {
	memDateCache.Flush()
}

func ClearCacheByUserID(userID uint) error {
	err := memDateCache.Delete(userID)
	if err != nil {
		return err
	}
	return nil
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

func FetchTodayMemDateByDeck(userID, deckID uint, skipCache bool) ([]models.MemDate, error) {
	db := database.DBConn // DB Conn
	t := time.Now()

	if !skipCache {
		// Check if cache exists
		if memDateCache.Exists(userID) {
			return memDateCache.Items(userID), nil
		}
	}
	var memDates []models.MemDate

	// Fetch today's memDates from DB
	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		userID).Joins("Card").Where("mem_dates.user_id = ? AND mem_dates.next_date < ? AND accesses.permission >= ? AND accesses.toggle_today IS true AND mem_dates.deck_id = ?",
		userID, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour), models.AccessStudent, deckID).Order("next_date asc").Find(&memDates).Error; err != nil {
		return nil, err
	}
	// Cache memDates
	memDateCache.SetSlice(userID, memDates)

	return memDates, nil
}

func FetchTodayMemDate(userID uint) ([]models.MemDate, error) {
	db := database.DBConn // DB Conn
	t := time.Now()

	// Check if cache exists

	if memDateCache.Exists(userID) {
		return memDateCache.Items(userID), nil
	}

	var memDates []models.MemDate

	// Fetch today's memDates from DB
	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		userID).Joins("Card").Where("mem_dates.user_id = ? AND mem_dates.next_date < ? AND accesses.permission >= ? AND accesses.toggle_today IS true",
		userID, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour), models.AccessStudent).Order("next_date asc").Find(&memDates).Error; err != nil {
		return nil, err
	}

	// Cache memDates
	memDateCache.SetSlice(userID, memDates)

	return memDates, nil
}

func GenerateResponseCardMap(memDates []models.MemDate, userID uint) (map[uint][]models.ResponseCard, error) {
	m := make(map[uint][]models.ResponseCard)

	wg := new(sync.WaitGroup)
	responseCard := new(models.ResponseCard)

	M := 12 // Number of handle per goroutine. Benchmark this value to optimize performance. (12 has been doing well)

	var workers int

	if len(memDates) < M {
		workers = 1
		M = len(memDates)
	} else {
		workers = len(memDates) / M
	}

	wg.Add(workers)

	ch := make(chan models.ResponseCard, len(memDates))

	var mutex = &sync.Mutex{}

	for i := 0; i < workers; i++ {
		hi, lo := i*M, (i+1)*M
		if i == workers-1 {
			lo = len(memDates)
		}

		subMemDates := memDates[hi:lo]
		go func() {
			for index := range subMemDates {
				mutex.Lock()
				answersList := GenerateMCQ(&subMemDates[index], userID)
				responseCard.Set(&subMemDates[index], answersList)
				ch <- *responseCard
				mutex.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)

	for toto := range ch {
		m[toto.Card.DeckID] = append(m[toto.Card.DeckID], toto)
	}

	return m, nil
}
