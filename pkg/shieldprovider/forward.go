package shieldprovider

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

// ForwardRootChain moves the root chain forward using HKDF based on 512 bit BLAKE2b hash function.
func (provider *Provider) ForwardRootChain(
	rootKey *keys.Root,
	sharedSecretKey *keys.SharedSecret,
) (*keys.Root, *keys.Master, *keys.Header, error) {
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
	masterKey := keys.NewMaster(output[32:64])
	nextHeaderKey := keys.NewHeader(output[64:96])

	return newRootKey, masterKey, nextHeaderKey, nil
}
