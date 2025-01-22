package shieldprovider

import (
	"crypto/ecdh"
	"crypto/rand"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

// GeneratePrivateKey generates a new X25519 private key using crypto/rand package.
func (provider Provider) GeneratePrivateKey() (*keys.Private, error) {
	foreignKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	key := keys.NewPrivate(foreignKey.Bytes())

	return key, nil
}
