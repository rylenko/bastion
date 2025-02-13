package messagechaincommon

import (
	"crypto/hmac"
	"fmt"
	"hash"

	"golang.org/x/crypto/blake2b"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

type Crypto struct{}

func NewCrypto() Crypto {
	return Crypto{}
}

func (c Crypto) AdvanceChain(masterKey keys.MessageMaster) (keys.MessageMaster, keys.Message, error) {
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
