package infrastructures_test

import (
	"testing"

	"github.com/memnix/memnix-rest/infrastructures"
)

func TestGetDBConnInstance(t *testing.T) {
	// Call GetDBConnInstance
	dbInstance := infrastructures.GetDBConnInstance()

	// Check if the returned instance is not nil
	if dbInstance == nil {
		t.Error("Failed to get DB connection instance")
	}

	// Call GetDBConnInstance again
	dbInstance2 := infrastructures.GetDBConnInstance()

	// Check if the second instance is the same as the first instance
	if dbInstance != dbInstance2 {
		t.Error("DB connection instances are not the same")
	}
}
