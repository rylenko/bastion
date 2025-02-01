package receivingchain

import (
	"crypto/hmac"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type Crypto interface {
	AdvanceChain(masterKey *keys.MessageMaster) (*keys.MessageMaster, *keys.Message, error)
}

type crypto struct{}

func newCrypto() Crypto {
	return crypto{}
}

func (crypto crypto) AdvanceChain(masterKey *keys.MessageMaster) (*keys.MessageMaster, *keys.Message, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrNewHash, err)
	}

	if masterKey == nil {
		return nil, nil, fmt.Errorf("%w: message master key is nil", ErrInvalidValue)
	}

	mac := hmac.New(func() hash.Hash { return hasher }, masterKey.Bytes())

	const masterKeyByte = 0x02
	if _, err := mac.Write([]byte{masterKeyByte}); err != nil {
		return nil, nil, fmt.Errorf("%w: message master key: %w", ErrMAC, err)
	}

	newMasterKey := keys.NewMessageMaster(mac.Sum(nil))
	mac.Reset()

	const messageKeyByte = 0x01
	if _, err := mac.Write([]byte{messageKeyByte}); err != nil {
		return nil, nil, fmt.Errorf("%w: message key: %w", ErrMAC, err)
	}

	messageKey := keys.NewMessage(mac.Sum(nil))

	return newMasterKey, messageKey, nil
}
