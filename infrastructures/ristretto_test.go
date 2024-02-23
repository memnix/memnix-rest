package infrastructures_test

import (
	"testing"

	"github.com/memnix/memnix-rest/infrastructures"
)

func TestGetCacheInstance(t *testing.T) {
	cacheInstance := infrastructures.GetCacheInstance()

	// Assert that the returned cache instance is not nil
	if cacheInstance == nil {
		t.Errorf("Expected cache instance to not be nil, but got nil")
	}
}

func TestCreateRistrettoInstance(t *testing.T) {
	// Create a mock RistrettoConfig
	mockConfig := infrastructures.RistrettoConfig{
		NumCounters: 1000,
		MaxCost:     1000,
		BufferItems: 1000,
	}

	// Call the CreateRistrettoInstance function
	cacheInstance := infrastructures.CreateRistrettoInstance(mockConfig)

	// Assert that the returned cache instance is not nil
	if cacheInstance == nil {
		t.Errorf("Expected cache instance to not be nil, but got nil")
	}

	// Assert that the returned cache instance has the correct config
	if cacheInstance != nil && cacheInstance.Config != mockConfig {
		t.Errorf("Expected cache instance to have the correct config, but got %v", cacheInstance.Config)
	}
}

func TestGetRistrettoCache(t *testing.T) {
	// Create a mock RistrettoConfig
	mockConfig := infrastructures.RistrettoConfig{
		NumCounters: 1000,
		MaxCost:     1000,
		BufferItems: 1000,
	}

	// Call the CreateRistrettoInstance function
	cacheInstance := infrastructures.CreateRistrettoInstance(mockConfig)

	err := cacheInstance.CreateRistrettoCache()
	if err != nil {
		t.Errorf("Failed to create ristretto cache: %v", err)
	}

	// Assert that the returned cache instance is not nil
	if cacheInstance == nil {
		t.Errorf("Expected cache instance to not be nil, but got nil")
	}

	// Assert that the returned cache instance has the correct config
	if cacheInstance != nil && cacheInstance.Config != mockConfig {
		t.Errorf("Expected cache instance to have the correct config, but got %v", cacheInstance.Config)
	}

	// Call the GetRistrettoCache function
	ristrettoCache := infrastructures.GetRistrettoCache()

	// Assert that the returned ristretto cache is not nil
	if ristrettoCache == nil {
		t.Errorf("Expected ristretto cache to not be nil, but got nil")
	}
}
