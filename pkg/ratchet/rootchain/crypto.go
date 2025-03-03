package rootchain

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

const kdfOutputLen = 3 * 32

var cryptoAdvanceChainKDFInfo = []byte("advance root chain")

type Crypto interface {
	AdvanceChain(rootKey keys.Root, sharedKey keys.Shared) (keys.Root, keys.MessageMaster, keys.Header, error)
}

type crypto struct{}

func newCrypto() Crypto {
	return crypto{}
}

func (crypto crypto) AdvanceChain(
	rootKey keys.Root,
	sharedKey keys.Shared,
) (keys.Root, keys.MessageMaster, keys.Header, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return keys.Root{}, keys.MessageMaster{}, keys.Header{}, fmt.Errorf("new hash: %w", err)
	}

	kdf := hkdf.New(func() hash.Hash { return hasher }, sharedKey.Bytes, rootKey.Bytes, cryptoAdvanceChainKDFInfo)

	output := make([]byte, kdfOutputLen)
	if _, err := io.ReadFull(kdf, output); err != nil {
		return keys.Root{}, keys.MessageMaster{}, keys.Header{}, fmt.Errorf("KDF: %w", err)
	}

	newRootKey := keys.Root{Bytes: output[:32]}
	messageMasterKey := keys.MessageMaster{Bytes: output[32:64]}
	nextHeaderKey := keys.Header{Bytes: output[64:]}

	return newRootKey, messageMasterKey, nextHeaderKey, nil
}
