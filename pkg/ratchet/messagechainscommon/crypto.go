package messagechainscommon

import (
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

const cryptoMessageCipherKDFOutputLen = cipher.KeySize + cipher.NonceSizeX

var (
	cryptoMessageCipherKDFSalt = make([]byte, cryptoMessageCipherKDFOutputLen)

	cryptoMessageCipherKDFInfo = []byte("message cipher")
)

func DeriveMessageCipherKeyAndNonce(messageKey keys.Message) ([]byte, []byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new hash: %w", err)
	}

	kdf := hkdf.New(
		func() hash.Hash { return hasher }, messageKey.Bytes, cryptoMessageCipherKDFSalt, cryptoMessageCipherKDFInfo)

	output := make([]byte, cryptoMessageCipherKDFOutputLen)
	if _, err := io.ReadFull(kdf, output); err != nil {
		return nil, nil, fmt.Errorf("KDF: %w", err)
	}

	return output[:cipher.KeySize], output[cipher.KeySize : cipher.KeySize+cipher.NonceSizeX], nil
}
