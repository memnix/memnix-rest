package crypto

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/pkg/errors"
)

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error generating keys")
	}

	return publicKey, privateKey, nil
}
