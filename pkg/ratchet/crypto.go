package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

var (
	cryptoEncryptHKDFInfo = []byte("encrypt participant data")
	cryptoEncryptHKDFSalt = make([]byte, blake2b.Size)
)

type Crypto interface {
	ComputeSharedSecretKey(privateKey *keys.Private, publicKey *keys.Public) (*keys.SharedSecret, error)
	Encrypt(key *keys.Message, data, auth []byte) ([]byte, error)
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
		return nil, fmt.Errorf("map to foreign private key: %w", err)
	}

	if publicKey == nil {
		return nil, fmt.Errorf("%w: public key is nil", ErrInvalidValue)
	}

	foreignPublicKey, err := c.curve.NewPublicKey(publicKey.Bytes())
	if err != nil {
		return nil, fmt.Errorf("map to foreign public key: %w", err)
	}

	sharedSecretKeyBytes, err := foreignPrivateKey.ECDH(foreignPublicKey)
	if err != nil {
		return nil, fmt.Errorf("Diffie-Hellman: %w", err)
	}

	sharedSecretKey := keys.NewSharedSecret(sharedSecretKeyBytes)

	return sharedSecretKey, nil
}

func (c *crypto) Encrypt(key *keys.Message, data, auth []byte) ([]byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, fmt.Errorf("new hash: %w", err)
	}

	hkdf := hkdf.New(func() hash.Hash { return hasher }, key.Bytes(), cryptoEncryptHKDFSalt, cryptoEncryptHKDFInfo)

	hkdfOutput := make([]byte, chacha20poly1305.KeySize+chacha20poly1305.NonceSizeX)
	if _, err := io.ReadFull(hkdf, hkdfOutput); err != nil {
		return nil, fmt.Errorf("KDF: %w", err)
	}

	cipherKey := hkdfOutput[:chacha20poly1305.KeySize]
	cipherNonce := hkdfOutput[len(hkdfOutput)-chacha20poly1305.NonceSizeX:]

	cipher, err := chacha20poly1305.NewX(cipherKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	encryptedData := cipher.Seal(nil, cipherNonce, data, auth)

	return encryptedData, nil
}

func (c *crypto) GeneratePrivateKey() (*keys.Private, error) {
	foreignKey, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	key := keys.NewPrivate(foreignKey.Bytes())

	return key, nil
}
