package queries

import (
	"fmt"
	"github.com/memnix/memnixrest/pkg/database"
	"github.com/memnix/memnixrest/pkg/models"
	"sync"
	"time"
)

var cache = make(map[uint]map[uint]models.MemDate)

func exportCache(userID uint) []models.MemDate {
	fmt.Println("export")
	memDates := make([]models.MemDate, 0, len(cache[userID]))
	for _, v := range cache[userID] {
		memDates = append(memDates, v)
	}
	return memDates
}

func importCache(userID uint, memDates []models.MemDate) {
	fmt.Println("import")
	cache[userID] = make(map[uint]models.MemDate, len(memDates))
	for _, v := range memDates {
		cache[userID][v.CardID] = v
	}
}

func FetchTodayMemDate(userID uint) ([]models.MemDate, error) {
	db := database.DBConn // DB Conn
	t := time.Now()

	var memDates []models.MemDate

	if cache[userID] != nil {
		fmt.Println("cache")
		return exportCache(userID), nil
	}

	if err := db.Joins(
		"left join accesses ON mem_dates.deck_id = accesses.deck_id AND accesses.user_id = ?",
		userID).Joins("Card").Where("mem_dates.user_id = ? AND mem_dates.next_date < ? AND accesses.permission >= ? AND accesses.toggle_today IS true",
		userID, t.AddDate(0, 0, 1).Add(
			time.Duration(-t.Hour())*time.Hour), models.AccessStudent).Order("next_date asc").Find(&memDates).Error; err != nil {
		return nil, err
	}

	importCache(userID, memDates)

	return memDates, nil
}

func GenerateResponseCardMap(memDates []models.MemDate, userID uint) (map[uint][]models.ResponseCard, error) {
	m := make(map[uint][]models.ResponseCard)

	wg := new(sync.WaitGroup)
	responseCard := new(models.ResponseCard)

	workers := 10

	if len(memDates) < workers {
		workers = 1
	} else if len(memDates) < workers*2 {
		workers = 4
	}

	M := len(memDates) / workers

	wg.Add(workers)

	ch := make(chan models.ResponseCard, len(memDates))

	for i := 0; i < workers; i++ {
		hi, lo := i*M, (i+1)*M
		if i == workers-1 {
			lo = len(memDates)
		}

		subMemDates := memDates[hi:lo]
		go func() {
			for index := range subMemDates {
				answersList := GenerateMCQ(&subMemDates[index], userID)
				responseCard.Set(&subMemDates[index], answersList)
				ch <- *responseCard
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
