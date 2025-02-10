package ratchet

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"unsafe"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/errors"
	"github.com/rylenko/bastion/pkg/ratchet/header"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

var (
	cryptoEncryptHeaderHKDFInfo = []byte("encrypt participant header")
	cryptoEncryptHeaderHKDFSalt = make([]byte, blake2b.Size)
	cryptoEncryptHKDFInfo       = []byte("encrypt participant data")
	cryptoEncryptHKDFSalt       = cryptoEncryptHeaderHKDFSalt

	uint64Size = int(unsafe.Sizeof(uint64(0)))
)

type Crypto interface {
	ComputeSharedSecretKey(privateKey *keys.Private, publicKey *keys.Public) (*keys.SharedSecret, error)
	Encrypt(key *keys.Message, data, auth []byte) ([]byte, error)
	EncryptHeader(key *keys.Header, header *header.Header) ([]byte, error)
	GenerateKeyPair() (*keys.Private, *keys.Public, error)
}

type crypto struct {
	curve ecdh.Curve
}

func newCrypto() Crypto {
	return &crypto{curve: ecdh.X25519()}
}

func (c crypto) ComputeSharedSecretKey(privateKey *keys.Private, publicKey *keys.Public) (*keys.SharedSecret, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("%w: private key is nil", errors.ErrInvalidValue)
	}

	foreignPrivateKey, err := c.curve.NewPrivateKey(privateKey.Bytes)
	if err != nil {
		return nil, fmt.Errorf("map to foreign private key: %w", err)
	}

	if publicKey == nil {
		return nil, fmt.Errorf("%w: public key is nil", errors.ErrInvalidValue)
	}

	foreignPublicKey, err := c.curve.NewPublicKey(publicKey.Bytes)
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

func (c crypto) Encrypt(key *keys.Message, data, auth []byte) ([]byte, error) {
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

	encryptedData := cipher.Seal(nil, cipherNonce, data, auth)

	return encryptedData, nil
}

func (c crypto) EncryptHeader(key *keys.Header, header *header.Header) ([]byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, fmt.Errorf("new hash: %w", err)
	}

	if key == nil {
		return nil, fmt.Errorf("%w: key is nil", errors.ErrInvalidValue)
	}

	// The header encryption key is repeated for each message until the next ratchet. Therefore, it is important not to
	// repeat the nonce when encrypting. Here, the HKDF is used based on the message number and the accepted header key.
	// This pair is unique, since the header key is unique, and after resetting the message number for the next sending
	// chain, we will get a different header key.
	hkdfKey := make([]byte, uint64Size+len(key.Bytes))
	binary.LittleEndian.PutUint64(hkdfKey[:uint64Size], header.MessageNumber)
	copy(hkdfKey[uint64Size:], key.Bytes)

	hkdf := hkdf.New(func() hash.Hash { return hasher }, hkdfKey, cryptoEncryptHeaderHKDFSalt, cryptoEncryptHeaderHKDFInfo)

	hkdfOutput := make([]byte, chacha20poly1305.KeySize+chacha20poly1305.NonceSizeX)
	if _, err := io.ReadFull(hkdf, hkdfOutput); err != nil {
		return nil, fmt.Errorf("KDF: %w", err)
	}

	cipherKey := hkdfOutput[:chacha20poly1305.KeySize]
	cipherNonce := hkdfOutput[len(hkdfOutput)-chacha20poly1305.NonceSizeX:]

	if header == nil {
		return nil, fmt.Errorf("%w: header is nil", errors.ErrInvalidValue)
	}

	headerBytes, err := header.Encode()
	if err != nil {
		return nil, fmt.Errorf("encode header: %w", err)
	}

	cipher, err := chacha20poly1305.NewX(cipherKey)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	encryptedData := cipher.Seal(nil, cipherNonce, headerBytes, nil)

	return encryptedData, nil
}

func (c crypto) GenerateKeyPair() (*keys.Private, *keys.Public, error) {
	foreignPrivateKey, err := c.curve.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	privateKey := keys.NewPrivate(foreignPrivateKey.Bytes())
	publicKey := keys.NewPublic(foreignPrivateKey.PublicKey().Bytes())

	return privateKey, publicKey, nil
}
