package sendingchain

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
	EncryptHeader(key keys.Header, header header.Header) ([]byte, error)
}

type crypto struct{}

func newCrypto() Crypto {
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

func (c crypto) EncryptHeader(headerKey keys.Header, header header.Header) ([]byte, error) {
	key, nonce, err := messagechainscommon.DeriveHeaderCipherKeyAndNonce(headerKey, header.MessageNumber)
	if err != nil {
		return nil, fmt.Errorf("derive key and nonce: %w", err)
	}

	cipher, err := cipher.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	return cipher.Seal(nil, nonce, header.Encode(), nil), nil
}
