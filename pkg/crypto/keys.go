package crypto

import (
	"log"
	"sync"

	"golang.org/x/crypto/ed25519"
)

type KeyManager struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

var (
	keyManagerInstance *KeyManager //nolint:gochecknoglobals //Singleton
	keyManagerOnce     sync.Once   //nolint:gochecknoglobals //Singleton
)

func GetKeyManagerInstance() *KeyManager {
	keyManagerOnce.Do(func() {
		keyManagerInstance = &KeyManager{}
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

	log.Println("✅ Created ed25519 keys")

	return nil
}
