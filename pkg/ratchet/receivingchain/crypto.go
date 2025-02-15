package receivingchain

import (
	"crypto/hmac"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"

	"github.com/rylenko/bastion/pkg/ratchet/header"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/messagechainscommon"
)

type Crypto interface {
	AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error)
	DecryptHeader(key keys.Header, encryptedHeader []byte, messageNumber uint64) (header.Header, error)
	DecryptMessage(key keys.Message, encryptedData, auth []byte) ([]byte, error)
}

type crypto struct{}

func newCrypto() crypto {
	return crypto{}
}

func (c crypto) AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return keys.MessageMaster{}, keys.Message{}, fmt.Errorf("new hash: %w", err)
	}

	mac := hmac.New(func() hash.Hash { return hasher }, masterKey.Bytes)

	const masterKeyByte = 0x02
	if _, err := mac.Write([]byte{masterKeyByte}); err != nil {
		return keys.MessageMaster{}, keys.Message{}, fmt.Errorf("write %d byte to MAC: %w", masterKeyByte, err)
	}

	newMasterKey := keys.MessageMaster{Bytes: mac.Sum(nil)}
	mac.Reset()

	const messageKeyByte = 0x01
	if _, err := mac.Write([]byte{messageKeyByte}); err != nil {
		return keys.MessageMaster{}, keys.Message{}, fmt.Errorf("write %d byte to MAC: %w", messageKeyByte, err)
	}

	messageKey := keys.Message{Bytes: mac.Sum(nil)}

	return newMasterKey, messageKey, nil
}

func (c crypto) DecryptHeader(key keys.Header, encryptedHeader []byte, messageNumber uint64) (header.Header, error) {
	cipherKey, nonce, err := messagechainscommon.DeriveHeaderCipherKeyAndNonce(key, messageNumber)
	if err != nil {
		return header.Header{}, fmt.Errorf("derive key and nonce: %w", err)
	}

	headerBytes, err := c.decrypt(cipherKey, nonce, encryptedHeader, nil)
	if err != nil {
		return header.Header{}, err
	}

	decryptedHeader, err := header.Decode(headerBytes)
	if err != nil {
		return header.Header{}, fmt.Errorf("decode decrypted header: %w", err)
	}

	return decryptedHeader, nil
}

func (c crypto) DecryptMessage(key keys.Message, encryptedData, auth []byte) ([]byte, error) {
	cipherKey, nonce, err := messagechainscommon.DeriveMessageCipherKeyAndNonce(key)
	if err != nil {
		return nil, fmt.Errorf("derive key and nonce: %w", err)
	}

	return c.decrypt(cipherKey, nonce, encryptedData, auth)
}

func (c crypto) decrypt(key, nonce, encryptedData, auth []byte) ([]byte, error) {
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
