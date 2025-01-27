package rootchain

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/sapphire/pkg/shield/keys"
)

type Crypto interface {
	AdvanceChain(
		rootKey *keys.Root,
		sharedSecretKey *keys.SharedSecret,
	) (*keys.Root, *keys.MessageMaster, *keys.Header, error)
}

type crypto struct{}

func newCrypto() Crypto {
	return crypto{}
}

func (crypto crypto) AdvanceChain(
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
