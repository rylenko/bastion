package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

const cryptoCipherKDFOutputLen = cipher.KeySize + cipher.NonceSizeX

var (
	cryptoCipherKDFInfo = []byte("ratchet")
	cryptoCipherKDFSalt = make([]byte, cryptoCipherKDFOutputLen)
)

type Crypto interface {
	ComputeSharedKey(privateKey keys.Private, publicKey keys.Public) (keys.Shared, error)
	Decrypt(key keys.Message, encryptedData, auth []byte) ([]byte, error)
	Encrypt(key keys.Message, data, auth []byte) ([]byte, error)
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

func (c crypto) Decrypt(messageKey keys.Message, encryptedData, auth []byte) ([]byte, error) {
	key, nonce, err := deriveCipherKeyAndNonce(messageKey)
	if err != nil {
		return nil, fmt.Errorf("derive key and nonce: %w", err)
	}

	cipher, err := cipher.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	data, err := cipher.Open(nil, nonce, encryptedData, auth)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return data, nil
}

func (c crypto) Encrypt(messageKey keys.Message, data, auth []byte) ([]byte, error) {
	key, nonce, err := deriveCipherKeyAndNonce(messageKey)
	if err != nil {
		return nil, fmt.Errorf("derive key and nonce: %w", err)
	}

	cipher, err := cipher.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	return cipher.Seal(nil, nonce, data, auth), nil
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

func deriveCipherKeyAndNonce(messageKey keys.Message) ([]byte, []byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new hash: %w", err)
	}

	hkdf := hkdf.New(func() hash.Hash { return hasher }, messageKey.Bytes, cryptoCipherKDFSalt, cryptoCipherKDFInfo)

	output := make([]byte, cryptoCipherKDFOutputLen)
	if _, err := io.ReadFull(hkdf, output); err != nil {
		return nil, nil, fmt.Errorf("KDF: %w", err)
	}

	return output[:cipher.KeySize], output[cipher.KeySize : cipher.KeySize+cipher.NonceSizeX], nil
}
