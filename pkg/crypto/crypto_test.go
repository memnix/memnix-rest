package crypto

import (
	"testing"
)

// MockCrypto is a mock implementation of the ICrypto interface for testing.
type MockCrypto struct{}

func (mc MockCrypto) Hash(_ string) ([]byte, error) {
	// Implement a mock hashing function for testing
	return []byte("mocked-hash"), nil
}

func (mc MockCrypto) Verify(_ string, hash []byte) (bool, error) {
	// Implement a mock verify function for testing
	return string(hash) == "mocked-hash", nil
}

func TestCrypto_Hash(t *testing.T) {
	mockCrypto := MockCrypto{}
	cryptoHelper = Crypto{crypto: mockCrypto}

	password := "my_password"
	hashedPassword, err := cryptoHelper.Hash(password)
	if err != nil {
		t.Errorf("Hash() returned an error: %v", err)
	}

	if string(hashedPassword) != "mocked-hash" {
		t.Errorf("Hash() did not return the expected hashed value")
	}
}

func TestCrypto_Verify(t *testing.T) {
	mockCrypto := MockCrypto{}
	cryptoHelper = Crypto{crypto: mockCrypto}

	password := "my_password"
	hashedPassword, _ := cryptoHelper.Hash(password)

	match, err := cryptoHelper.Verify(password, hashedPassword)
	if err != nil {
		t.Errorf("Verify() returned an error: %v", err)
	}

	if !match {
		t.Errorf("Verify() did not return a match for the hashed password")
	}
}
