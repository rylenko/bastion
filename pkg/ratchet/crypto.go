package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type Crypto interface {
	ComputeSharedKey(privateKey keys.Private, publicKey keys.Public) (keys.Shared, error)
	GenerateKeyPair() (keys.Private, keys.Public, error)
}

type crypto struct {
	curve ecdh.Curve
}

func newCrypto() Crypto {
	return crypto{curve: ecdh.X25519()}
}

func (c crypto) ComputeSharedKey(privateKey keys.Private, publicKey keys.Public) (keys.Shared, error) {
	foreignPrivateKey, err := c.curve.NewPrivateKey(privateKey.Bytes)
	if err != nil {
		return keys.Shared{}, fmt.Errorf("map to foreign private key: %w", err)
	}

	foreignPublicKey, err := c.curve.NewPublicKey(publicKey.Bytes)
	if err != nil {
		return keys.Shared{}, fmt.Errorf("map to foreign public key: %w", err)
	}

	sharedKeyBytes, err := foreignPrivateKey.ECDH(foreignPublicKey)
	if err != nil {
		return keys.Shared{}, fmt.Errorf("Diffie-Hellman: %w", err)
	}

	return keys.Shared{Bytes: sharedKeyBytes}, nil
}

func (c crypto) GenerateKeyPair() (keys.Private, keys.Public, error) {
	foreignPrivateKey, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return keys.Private{}, keys.Public{}, err
	}

	privateKey := keys.Private{Bytes: foreignPrivateKey.Bytes()}
	publicKey := keys.Public{Bytes: foreignPrivateKey.PublicKey().Bytes()}

	return privateKey, publicKey, nil
}
