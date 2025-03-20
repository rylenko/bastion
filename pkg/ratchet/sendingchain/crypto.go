package sendingchain

import (
	"crypto/hmac"
	"crypto/rand"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"

	"github.com/rylenko/bastion/pkg/ratchet/header"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/messagechainscommon"
	"github.com/rylenko/bastion/pkg/utils"
)

type Crypto interface {
	AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error)
	EncryptHeader(key keys.Header, header header.Header) ([]byte, error)
	EncryptMessage(key keys.Message, data, auth []byte) ([]byte, error)
}

type defaultCrypto struct{}

func newDefaultCrypto() Crypto {
	return defaultCrypto{}
}

func (c defaultCrypto) AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error) {
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

func (c defaultCrypto) EncryptHeader(key keys.Header, header header.Header) ([]byte, error) {
	var nonce [cipher.NonceSizeX]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("generate random nonce: %w", err)
	}

	encryptedHeader, err := c.encrypt(key.Bytes, nonce[:], header.Encode(), nil)
	if err != nil {
		return nil, err
	}

	return utils.ConcatByteSlices(nonce[:], encryptedHeader), nil
}

func (c defaultCrypto) EncryptMessage(key keys.Message, data, auth []byte) ([]byte, error) {
	cipherKey, nonce, err := messagechainscommon.DeriveMessageCipherKeyAndNonce(key)
	if err != nil {
		return nil, fmt.Errorf("derive key and nonce: %w", err)
	}

	return c.encrypt(cipherKey, nonce, data, auth)
}

func (c defaultCrypto) encrypt(key, nonce, data, auth []byte) ([]byte, error) {
	cipher, err := cipher.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	return cipher.Seal(nil, nonce, data, auth), nil
}
