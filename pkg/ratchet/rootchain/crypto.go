package rootchain

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

var cryptoAdvanceChainHKDFInfo = []byte("advance root chain")

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

	hkdf := hkdf.New(
		func() hash.Hash { return hasher }, sharedSecretKey.Bytes(), rootKey.Bytes(), cryptoAdvanceChainHKDFInfo)

	const hkdfOutputLen = 3 * 32
	hkdfOutput := make([]byte, hkdfOutputLen)

	if _, err := io.ReadFull(hkdf, hkdfOutput); err != nil {
		return nil, nil, nil, fmt.Errorf("%w: %w", ErrKDF, err)
	}

	newRootKey := keys.NewRoot(hkdfOutput[:32])
	messageMasterKey := keys.NewMessageMaster(hkdfOutput[32:64])
	nextHeaderKey := keys.NewHeader(hkdfOutput[64:96])

	return newRootKey, messageMasterKey, nextHeaderKey, nil
}
