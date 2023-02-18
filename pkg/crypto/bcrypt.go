package crypto

import (
	"github.com/corentings/kafejo/config"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// BcryptCrypto is the struct that holds the bcrypt crypto methods
type BcryptCrypto struct{}

// Hash hashes a password using the bcrypt algorithm
// password is the plaintext password to hash.
// Returns the hashed password, or an error on failure.
// The cost is set in the config file.
//
// see: https://godoc.org/golang.org/x/crypto/bcrypt
// see: utils/config.go for the default cost
func (*BcryptCrypto) Hash(password string) ([]byte, error) {
	key, err := bcrypt.GenerateFromPassword([]byte(password), config.BCryptCost)
	if err != nil {
		return []byte(""), err
	}
	return key, nil
}

// Verify compares a bcrypt hashed password with its possible plaintext equivalent.
// password is the plaintext password to verify.
// hash is the bcrypt hashed password.
// Returns nil on success, or an error on failure.
// Returns true if the password matches, false if it does not.
func (*BcryptCrypto) Verify(password string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return false, errors.Wrap(err, errors.New("error comparing bcrypt hash").Error())
	}
	return true, nil
}
