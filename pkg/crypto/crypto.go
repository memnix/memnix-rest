package crypto

// ICrypto is the interface for the crypto methods
// It's used to abstract the crypto methods used in the application
// so that they can be easily swapped out if needed
type ICrypto interface {
	// Hash hashes a password using the configured crypto method
	Hash(password string) ([]byte, error)
	// Verify compares a crypto hashed password with its possible plaintext equivalent
	Verify(password string, hash []byte) (bool, error)
}

// Crypto is the struct that holds the crypto methods
type Crypto struct {
	crypto ICrypto
}

// Hash hashes a password using the configured crypto method
// password is the plaintext password to hash.
// Returns the hashed password, or an error on failure.
func (c *Crypto) Hash(password string) ([]byte, error) {
	return c.crypto.Hash(password) //nolint:wrapcheck
}

// Verify compares a crypto hashed password with its possible plaintext equivalent
// password is the plaintext password to verify.
// hash is the bcrypt hashed password.
// Returns nil on success, or an error on failure.
// Returns true if the password matches, false if it does not.
func (c *Crypto) Verify(password string, hash []byte) (bool, error) {
	return c.crypto.Verify(password, hash) //nolint:wrapcheck
}
