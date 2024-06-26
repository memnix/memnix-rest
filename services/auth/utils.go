package auth

import (
	"context"

	"github.com/memnix/memnix-rest/pkg/crypto"
	"github.com/pkg/errors"
)

// GenerateEncryptedPassword generates a password hash using the crypto helper.
func GenerateEncryptedPassword(_ context.Context, password string) ([]byte, error) {
	hash, err := crypto.GetCryptoHelperInstance().GetCryptoHelper().Hash(password)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// ComparePasswords compares a hashed password with its possible plaintext equivalent.
//
// password is the plaintext password to verify.
// hash is the bcrypt hashed password.
//
// Returns true if the password matches, false if it does not.
// Returns nil on success, or an error on failure.
func ComparePasswords(_ context.Context, password string, hash []byte) (bool, error) {
	return crypto.GetCryptoHelperInstance().GetCryptoHelper().Verify(password, hash)
}

// VerifyPassword verifies a password
// Returns an error if the password is invalid.
func VerifyPassword(password string) error {
	// Convert password to byte array
	passwordBytes := []byte(password)
	if len(passwordBytes) < crypto.MinPasswordLength {
		return errors.New("password too short")
	}

	if len(passwordBytes) > crypto.MaxPasswordLength {
		return errors.New("password too long")
	}

	return nil
}
