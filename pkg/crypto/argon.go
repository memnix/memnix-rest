package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"github.com/corentings/kafejo/config"
	"github.com/pkg/errors"
	"golang.org/x/crypto/argon2"
	"strings"
)

// Argon2Crypto is the struct that holds the argon2id crypto methods
type Argon2Crypto struct{}

// Hash hashes a password using the argon2id algorithm
// password is the plaintext password to hash.
// Returns the hashed password, or an error on failure.
// The returned hash is in the format:
// $argon2id$v=VERSION$m=MEMORY,t=ITERATIONS,p=THREADS$SALT$KEY
//
// see: https://godoc.org/golang.org/x/crypto/argon2
func (*Argon2Crypto) Hash(password string) ([]byte, error) {
	cryptoConfig := config.PasswordConfig
	var key []byte

	salt, err := GenerateRandomBytes(cryptoConfig.SaltLen)
	if err != nil {
		return []byte(""), err
	}

	key = argon2.IDKey([]byte(password), salt, cryptoConfig.Iterations, cryptoConfig.Memory, cryptoConfig.Threads, cryptoConfig.KeyLen)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Key := base64.RawStdEncoding.EncodeToString(key)

	hash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, cryptoConfig.Memory, cryptoConfig.Iterations, cryptoConfig.Threads, b64Salt, b64Key)

	return []byte(hash), nil
}

// Verify compares a crypto hashed password with its possible plaintext equivalent
// password is the plaintext password to verify.
// hash is the argon2id hashed password.
// Returns nil on success, or an error on failure.
// Returns true if the password matches, false if it does not.
func (*Argon2Crypto) Verify(password string, hash []byte) (bool, error) {
	cryptoConfig, salt, key, err := DecodeHash(string(hash))
	if err != nil {
		return false, err
	}

	newKey := argon2.IDKey([]byte(password), salt, cryptoConfig.Iterations, cryptoConfig.Memory, cryptoConfig.Threads, cryptoConfig.KeyLen)

	return SecureCompare(key, newKey), nil
}

// GenerateRandomBytes generates a byte slice of the given length, filled with random bytes.
func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, errors.New("error generating random bytes")
	}

	return b, nil
}

// DecodeHash decodes a hash and returns the configuration, salt, and key.
// It's used for the argon2id algorithm
func DecodeHash(hash string) (*config.PasswordConfigStruct, []byte, []byte, error) {
	vals := strings.Split(hash, "$")
	if len(vals) != config.Argon2HashLength {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash format")
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version of argon2")
	}

	cryptoConfig := &config.PasswordConfigStruct{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &cryptoConfig.Memory, &cryptoConfig.Iterations, &cryptoConfig.Threads)
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash format")
	}

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash format")
	}
	cryptoConfig.SaltLen = uint32(len(salt))

	key, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, errors.New("invalid hash format")
	}
	cryptoConfig.KeyLen = uint32(len(key))

	return cryptoConfig, salt, key, nil
}

// SecureCompare returns true if the two byte slices are equal.
func SecureCompare(a, b []byte) bool {
	aLen := int32(len(a))
	bLen := int32(len(b))

	if subtle.ConstantTimeEq(aLen, bLen) == 0 {
		return false
	}
	if subtle.ConstantTimeCompare(a, b) == 1 {
		return true
	}

	return false
}
