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
	foreignPrivateKey, err := provider.mapToForeignPrivateKey(privateKey)
	if err != nil {
		return nil, fmt.Errorf("%w: private key: %w", ErrMapToForeignType, err)
	}

	foreignPublicKey, err := provider.mapToForeignPublicKey(publicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: public key: %w", ErrMapToForeignType, err)
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

	key, err := provider.mapFromForeignPrivateKey(foreignKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMapFromForeignType, err)
	}

	return key, nil
}
