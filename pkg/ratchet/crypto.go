package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/header"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/utils"
)

var (
	cryptoEncryptHeaderHKDFInfo = []byte("encrypt participant header")
	cryptoEncryptHeaderHKDFSalt = make([]byte, blake2b.Size)
	cryptoEncryptHKDFInfo       = []byte("encrypt participant data")
	cryptoEncryptHKDFSalt       = cryptoEncryptHeaderHKDFSalt
)

type Crypto interface {
	ComputeSharedKey(privateKey keys.Private, publicKey keys.Public) (keys.Shared, error)
	Encrypt(key keys.Message, data, auth []byte) ([]byte, error)
	EncryptHeader(key keys.Header, header header.Header) ([]byte, error)
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

func (c crypto) Encrypt(key keys.Message, data, auth []byte) ([]byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, fmt.Errorf("new hash: %w", err)
	}

	hkdf := hkdf.New(func() hash.Hash { return hasher }, key.Bytes, cryptoEncryptHKDFSalt, cryptoEncryptHKDFInfo)

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

	return cipher.Seal(nil, cipherNonce, data, auth), nil
}

func (c crypto) EncryptHeader(key keys.Header, header header.Header) ([]byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, fmt.Errorf("new hash: %w", err)
	}

	// The header encryption key is repeated for each message until the next ratchet. Therefore, it is important not to
	// repeat the nonce when encrypting. Here, the HKDF is used based on the message number and the accepted header key.
	// This pair is unique, since the header key is unique, and after resetting the message number for the next sending
	// chain, we will get a different header key.
	var messageNumberBytes [utils.Uint64Size]byte
	binary.LittleEndian.PutUint64(messageNumberBytes[:], header.MessageNumber)
	hkdfKey := utils.ConcatByteSlices(messageNumberBytes[:], key.Bytes)

	hkdf := hkdf.New(func() hash.Hash { return hasher }, hkdfKey, cryptoEncryptHeaderHKDFSalt, cryptoEncryptHeaderHKDFInfo)

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

	return cipher.Seal(nil, cipherNonce, header.Encode(), nil), nil
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
