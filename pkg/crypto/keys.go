package crypto

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ed25519"
)

type KeyManager struct {
	seed       io.Reader
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

var (
	keyManagerInstance *KeyManager //nolint:gochecknoglobals //Singleton
	keyManagerOnce     sync.Once   //nolint:gochecknoglobals //Singleton
)

func GetKeyManagerInstance() *KeyManager {
	keyManagerOnce.Do(func() {
		keyManagerInstance = &KeyManager{
			seed: rand.Reader,
		}
	})
	return keyManagerInstance
}

func (k *KeyManager) GetPrivateKey() ed25519.PrivateKey {
	return k.privateKey
}

func (k *KeyManager) GetPublicKey() ed25519.PublicKey {
	return k.publicKey
}

func (k *KeyManager) ParseEd25519Key() error {
	publicKey, privateKey, err := GenerateKeyPair()
	if err != nil {
		return err
	}

	k.privateKey = privateKey
	k.publicKey = publicKey

	return nil
}

func GenerateKeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(GetKeyManagerInstance().seed)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error generating keys")
	}

	return publicKey, privateKey, nil
}
