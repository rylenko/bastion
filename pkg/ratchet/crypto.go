package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type Crypto interface {
	ComputeSharedSecretKey(privateKey *keys.Private, publicKey *keys.Public) (*keys.SharedSecret, error)
	GeneratePrivateKey() (*keys.Private, error)
}

type crypto struct {
	curve ecdh.Curve
}

func newCrypto() Crypto {
	return &crypto{curve: ecdh.X25519()}
}

func (c *crypto) ComputeSharedSecretKey(privateKey *keys.Private, publicKey *keys.Public) (*keys.SharedSecret, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("%w: private key is nil", ErrInvalidValue)
	}

	foreignPrivateKey, err := c.curve.NewPrivateKey(privateKey.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: private key: %w", ErrForeignType, err)
	}

	if publicKey == nil {
		return nil, fmt.Errorf("%w: public key is nil", ErrInvalidValue)
	}

	foreignPublicKey, err := c.curve.NewPublicKey(publicKey.Bytes())
	if err != nil {
		return nil, fmt.Errorf("%w: public key: %w", ErrForeignType, err)
	}

	sharedSecretKeyBytes, err := foreignPrivateKey.ECDH(foreignPublicKey)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDiffieHellman, err)
	}

	sharedSecretKey := keys.NewSharedSecret(sharedSecretKeyBytes)

	return sharedSecretKey, nil
}

func (c *crypto) GeneratePrivateKey() (*keys.Private, error) {
	foreignKey, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	key := keys.NewPrivate(foreignKey.Bytes())

	return key, nil
}
