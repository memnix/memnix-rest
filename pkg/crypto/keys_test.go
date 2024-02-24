package crypto_test

import (
	"crypto/ed25519"
	"testing"

	"github.com/memnix/memnix-rest/pkg/crypto"
)

func TestGenerateKeyPair(t *testing.T) {
	publicKey, privateKey, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(publicKey) != ed25519.PublicKeySize {
		t.Errorf("Expected public key length to be %d, but got %d", ed25519.PublicKeySize, len(publicKey))
	}

	if len(privateKey) != ed25519.PrivateKeySize {
		t.Errorf("Expected private key length to be %d, but got %d", ed25519.PrivateKeySize, len(privateKey))
	}
}
