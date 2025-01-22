package shieldprovider

import (
	"crypto/rand"
	"fmt"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

// ComputeSharedSecretKey computes shared secret key using Diffie-Hellman algorithm and Curve25519.
func (provider *Provider) ComputeSharedSecretKey(
	privateKey *keys.Private,
	publicKey *keys.Public,
) (*keys.SharedSecret, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("%w: private key is nil", ErrInvalidValue)
	}

	foreignPrivateKey, err := provider.curve.NewPrivateKey(privateKey.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: private key: %w", ErrConvertToForeignType, err)
	}

	if publicKey == nil {
		return nil, fmt.Errorf("%w: public key is nil", ErrInvalidValue)
	}

	foreignPublicKey, err := provider.curve.NewPublicKey(publicKey.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: public key: %w", ErrConvertToForeignType, err)
	}

	sharedSecretKeyBytes, err := foreignPrivateKey.ECDH(foreignPublicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDiffieHellman, err)
	}

	sharedSecretKey := keys.NewSharedSecret(sharedSecretKeyBytes)

	return sharedSecretKey, nil
}

// GeneratePrivateKey generates a new X25519 private key using crypto/rand package.
func (provider *Provider) GeneratePrivateKey() (*keys.Private, error) {
	foreignKey, err := provider.curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	key := keys.NewPrivate(foreignKey.Bytes())

	return key, nil
}
