package infrastructures

import (
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/memnix/memnix-rest/config"
	"github.com/memnix/memnix-rest/pkg/env"
)

var (
	client    *meilisearch.Client
	searchKey *meilisearch.Key
)

func ConnectMeiliSearch(env env.IEnv) {
	var host string
	if config.IsDevelopment() {
		host = env.GetEnv("DEBUG_MEILISEARCH_HOST")
	} else {
		host = env.GetEnv("MEILISEARCH_HOST")
	}
	client = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   host,
		APIKey: env.GetEnv("MEILISEARCH_API_KEY"),
	})
}

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

func GetMeiliSearchClient() *meilisearch.Client {
	return client
}
