package shieldprovider

import (
	"crypto/hmac"
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

// ForwardMessageChain moves the message chain forward using HMAC based on 512 bit BLAKE2b hash function.
func (provider *Provider) ForwardMessageChain(
	messageMasterKey *keys.MessageMaster,
) (*keys.MessageMaster, *keys.Message, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %w", ErrNewHash, err)
	}

	if messageMasterKey == nil {
		return nil, nil, fmt.Errorf("%w: message master key is nil", ErrInvalidValue)
	}

	mac := hmac.New(func() hash.Hash { return hasher }, messageMasterKey.Bytes())

	const messageMasterKeyByte = 0x02
	if _, err := mac.Write([]byte{messageMasterKeyByte}); err != nil {
		return nil, nil, fmt.Errorf("%w: message master key: %w", ErrMAC, err)
	}

	newMessageMasterKey := keys.NewMessageMaster(mac.Sum(nil))
	mac.Reset()

	const messageKeyByte = 0x01
	if _, err := mac.Write([]byte{messageKeyByte}); err != nil {
		return nil, nil, fmt.Errorf("%w: message key: %w", ErrMAC, err)
	}

	messageKey := keys.NewMessage(mac.Sum(nil))

	return newMessageMasterKey, messageKey, nil
}

// ForwardRootChain moves the root chain forward using HKDF based on 512 bit BLAKE2b hash function.
func (provider *Provider) ForwardRootChain(
	rootKey *keys.Root,
	sharedSecretKey *keys.SharedSecret,
) (*keys.Root, *keys.MessageMaster, *keys.Header, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %w", ErrNewHash, err)
	}

	if sharedSecretKey == nil {
		return nil, nil, nil, fmt.Errorf("%w: shared secret key is nil", ErrInvalidValue)
	}

	if rootKey == nil {
		return nil, nil, nil, fmt.Errorf("%w: root key is nil", ErrInvalidValue)
	}

	const hkdfInfo = "ForwardRootChainHKDFInfo"
	hkdf := hkdf.New(func() hash.Hash { return hasher }, sharedSecretKey.Bytes(), rootKey.Bytes(), []byte(hkdfInfo))

	const outputLen = 3 * 32
	output := make([]byte, outputLen)

	if _, err := io.ReadFull(hkdf, output); err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %w", ErrKDF, err)
	}

	newRootKey := keys.NewRoot(output[:32])
	messageMasterKey := keys.NewMessageMaster(output[32:64])
	nextHeaderKey := keys.NewHeader(output[64:96])

	return newRootKey, messageMasterKey, nextHeaderKey, nil
}
