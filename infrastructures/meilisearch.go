package infrastructures

import (
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/env"
)

var (
	client    *meilisearch.Client // client is the MeiliSearch client.
	searchKey *meilisearch.Key    // searchKey is the MeiliSearch search key.
)

// ConnectMeiliSearch connects to MeiliSearch.
func ConnectMeiliSearch(env env.IEnv) {
	var host string
	var apiKey string
	if config.IsDevelopment() {
		host = env.GetEnv("DEBUG_MEILISEARCH_HOST")
		apiKey = env.GetEnv("DEBUG_MEILISEARCH_API_KEY")
	} else {
		host = env.GetEnv("MEILISEARCH_HOST")
		apiKey = env.GetEnv("MEILISEARCH_API_KEY")
	}
	client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   host,
		APIKey: apiKey,
	})
}

// CreateSearchKey creates a new search key.
func CreateSearchKey() error {
	key, err := client.CreateKey(&meilisearch.Key{
		Description: "Api search key",
		Actions:     []string{"search"},
		Indexes:     []string{"*"},
		ExpiresAt:   time.Time{},
	})
	if err != nil {
		searchKey = nil
		return err
	}
	searchKey = key
	return nil
}

// GetSearchKey returns the search key.
func GetSearchKey() (meilisearch.Key, error) {
	if searchKey != nil {
		return *searchKey, nil
	}
	err := CreateSearchKey()
	if err != nil {
		return meilisearch.Key{}, err
	}
	return *searchKey, nil
}

// GetMeiliSearchClient returns the MeiliSearch client.
func GetMeiliSearchClient() *meilisearch.Client {
	return client
}
