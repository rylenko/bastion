package messagechainscommon

import (
	"encoding/binary"
	"fmt"
	"hash"
	"io"

	"golang.org/x/crypto/blake2b"
	cipher "golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/rylenko/bastion/pkg/ratchet/keys"
	"github.com/rylenko/bastion/pkg/ratchet/utils"
)

const cryptoCipherKDFOutputLen = cipher.KeySize + cipher.NonceSizeX

var (
	cryptoCipherKDFSalt = make([]byte, cryptoCipherKDFOutputLen)

	cryptoHeaderCipherKDFInfo  = []byte("header cipher")
	cryptoMessageCipherKDFInfo = []byte("message cipher")
)

func DeriveMessageCipherKeyAndNonce(messageKey keys.Message) ([]byte, []byte, error) {
	return deriveCipherKeyAndNonce(messageKey.Bytes, cryptoMessageCipherKDFInfo)
}

func DeriveHeaderCipherKeyAndNonce(headerKey keys.Header, messageNumber uint64) ([]byte, []byte, error) {
	// The header encryption key is repeated for each message until the next ratchet. Therefore, it is important not to
	// repeat the nonce when encrypting. Here, the KDF is used based on the message number and the accepted header key.
	// This pair is unique, since the header key is unique, and after resetting the message number for the next sending
	// chain, we will get a different header key.
	var messageNumberBytes [utils.Uint64Size]byte
	binary.LittleEndian.PutUint64(messageNumberBytes[:], messageNumber)
	kdfKey := utils.ConcatByteSlices(messageNumberBytes[:], headerKey.Bytes)

	return deriveCipherKeyAndNonce(kdfKey, cryptoHeaderCipherKDFInfo)
}

func deriveCipherKeyAndNonce(kdfKey, info []byte) ([]byte, []byte, error) {
	hasher, err := blake2b.New512(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("new hash: %w", err)
	}

	kdf := hkdf.New(func() hash.Hash { return hasher }, kdfKey, cryptoCipherKDFSalt, info)

	output := make([]byte, cryptoCipherKDFOutputLen)
	if _, err := io.ReadFull(kdf, output); err != nil {
		return nil, nil, fmt.Errorf("KDF: %w", err)
	}

	return output[:cipher.KeySize], output[cipher.KeySize : cipher.KeySize+cipher.NonceSizeX], nil
}
